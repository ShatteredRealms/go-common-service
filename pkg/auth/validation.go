package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/ShatteredRealms/go-common-service/pkg/srospan"
	"github.com/WilSimpson/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrFailed         = errors.New("auth")
	ErrMissingGocloak = fmt.Errorf("%w: CA01", ErrFailed)
	ErrDoesNotExist   = fmt.Errorf("%w: CA02", ErrFailed)
)

func verifyClaims(ctx context.Context, client gocloak.GoCloakIface, realm string) (*jwt.Token, *SROClaims, error) {
	if client == nil {
		return nil, nil, ErrMissingGocloak
	}

	tokenString, err := extractToken(ctx)
	if err != nil {
		return nil, nil, err
	}

	var claims SROClaims
	token, err := client.DecodeAccessTokenCustomClaims(
		ctx,
		tokenString,
		realm,
		&claims,
	)

	if err != nil {
		log.Logger.WithContext(ctx).Infof("Error extracting claims: %v", err)
		return nil, nil, err
	}

	if !token.Valid {
		log.Logger.WithContext(ctx).Infof("Invalid token given from %s:%s", claims.Username, claims.Subject)
		return nil, nil, err
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.SourceOwnerId(claims.Subject),
		srospan.SourceOwnerUsername(claims.Username),
	)

	return token, &claims, nil
}
