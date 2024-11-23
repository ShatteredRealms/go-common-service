package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/metadata"
)

const (
	claimContextKey claimContextKeyType = iota
)

var (
	ErrMissingAuthorization = errors.New("missing authorization")
	ErrInvalidAuthorization = errors.New("invalid authorization scheme")
)

func RetrieveClaims(ctx context.Context) (claims *SROClaims, ok bool) {
	claims, ok = ctx.Value(claimContextKey).(*SROClaims)
	return
}

func insertClaims(ctx context.Context, claims *SROClaims) context.Context {
	return context.WithValue(ctx, claimContextKey, claims)
}

func extractToken(ctx context.Context) (string, error) {
	val := metautils.ExtractIncoming(ctx).Get(AuthorizationHeader)
	if val == "" {
		return "", ErrMissingAuthorization
	}

	if !strings.HasPrefix(val, AuthorizationScheme) {
		return "", ErrInvalidAuthorization
	}

	return val[len(AuthorizationScheme):], nil
}

func AddOutgoingToken(ctx context.Context, token string) context.Context {
	return addOutgoingAuthBearer(ctx, "Bearer "+token)
}

func PassOutgoing(ctx context.Context) context.Context {
	return addOutgoingAuthBearer(
		ctx,
		metautils.ExtractIncoming(ctx).Get("authorization"),
	)
}

func addOutgoingAuthBearer(ctx context.Context, token string) context.Context {
	md := metadata.New(
		map[string]string{
			"authorization": token,
		},
	)

	return metadata.AppendToOutgoingContext(metadata.NewOutgoingContext(ctx, md))
}
