package auth

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/log"
	gocloak "github.com/WilSimpson/gocloak/v13"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnauthorized = status.New(codes.Unauthenticated, "not authorized")
)

type claimContextKeyType int8

var (
	publicMethods = make(map[string]struct{}, 10)
)

func init() {
	publicMethods["/sro.HealthService/Health"] = struct{}{}
}

func RegisterPublicServiceMethods(methods ...string) {
	for _, method := range methods {
		publicMethods[method] = struct{}{}
	}
}

func AuthFunc(kcClient gocloak.GoCloakIface, realm string) auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		_, claims, err := verifyClaims(ctx, kcClient, realm)

		if err != nil {
			log.Logger.WithContext(ctx).Infof("Verifying claims failed: %s", err)
			return nil, ErrUnauthorized.Err()
		}

		return insertClaims(ctx, claims), nil
	}
}

func NotPublicServiceMatcher(ctx context.Context, callMeta interceptors.CallMeta) bool {
	_, ok := publicMethods[callMeta.FullMethod()]
	log.Logger.WithContext(ctx).Debugf("Verify Auth (%t): %s", !ok, callMeta.FullMethod())
	return !ok
}
