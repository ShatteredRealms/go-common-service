package srv

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/WilSimpson/gocloak/v13"
	"github.com/golang-jwt/jwt/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	parser = jwt.Parser{}
)

type Context struct {
	Config         *config.BaseConfig
	KeycloakClient *gocloak.GoCloak
	Tracer         trace.Tracer

	jwt            *gocloak.JWT
	tokenExpiresAt time.Time
}

func NewContext(config *config.BaseConfig, service string) *Context {
	srvCtx := &Context{
		Config:         config,
		KeycloakClient: gocloak.NewClient(config.Keycloak.BaseURL),
		Tracer:         otel.Tracer(service),
	}

	srvCtx.KeycloakClient.RegisterMiddlewares(gocloak.OpenTelemetryMiddleware)
	log.Logger.Level = config.LogLevel

	return srvCtx
}

func (srvCtx *Context) loginClient(ctx context.Context) (*gocloak.JWT, error) {
	var err error
	srvCtx.jwt, err = srvCtx.KeycloakClient.LoginClient(
		ctx,
		srvCtx.Config.Keycloak.ClientId,
		srvCtx.Config.Keycloak.ClientSecret,
		srvCtx.Config.Keycloak.Realm,
	)
	if err != nil {
		return nil, fmt.Errorf("login keycloak: %v", err)
	}

	claims := &jwt.RegisteredClaims{}
	_, _, err = parser.ParseUnverified(srvCtx.jwt.AccessToken, claims)
	if err != nil {
		log.Logger.Errorf("parsing access token: %v", err)
		return srvCtx.jwt, nil
	}

	// Remove 5 seconds to ensure there are no race cases with expiration
	srvCtx.tokenExpiresAt = claims.ExpiresAt.Time.Add(-5 * time.Second)
	return srvCtx.jwt, nil
}

func (srvCtx *Context) GetJWT(ctx context.Context) (*gocloak.JWT, error) {
	if srvCtx.jwt != nil && time.Now().Before(srvCtx.tokenExpiresAt) {
		return srvCtx.jwt, nil
	}

	return srvCtx.loginClient(ctx)
}

func (srvCtx *Context) CreateRoles(ctx context.Context, roles *[]*gocloak.Role) error {
	ctx, span := srvCtx.Tracer.Start(ctx, "roles.create")
	defer span.End()
	jwtToken, err := srvCtx.GetJWT(ctx)
	if err != nil {
		return err
	}

	var errs error
	for _, role := range *roles {
		_, err := srvCtx.KeycloakClient.CreateClientRole(
			ctx,
			jwtToken.AccessToken,
			srvCtx.Config.Keycloak.Realm,
			srvCtx.Config.Keycloak.Id,
			*role,
		)

		// Code 409 is conflict
		if err != nil {
			if err.(*gocloak.APIError).Code != 409 {
				span.SetAttributes(attribute.String("role."+*role.Name, "error"))
				errs = errors.Join(errs, fmt.Errorf("creating role %s: %v", *role.Name, err))
			} else {
				span.SetAttributes(attribute.String("role."+*role.Name, "exists"))
			}
		} else {
			span.SetAttributes(attribute.String("role."+*role.Name, "created"))
		}
	}

	return errs
}
