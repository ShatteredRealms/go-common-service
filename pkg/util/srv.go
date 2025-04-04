package util

import (
	"context"
	"fmt"

	sroauth "github.com/ShatteredRealms/go-common-service/pkg/auth"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	gocloak "github.com/WilSimpson/gocloak/v13"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GrpcDialOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

func GrpcClientWithOtel(address string) (*grpc.ClientConn, error) {
	return grpc.NewClient(address, GrpcDialOpts()...)
}

func InitServerDefaults(kcClient gocloak.GoCloakIface, realm string) (*grpc.Server, *runtime.ServeMux) {
	opts := []logging.Option{
		logging.WithCodes(logging.DefaultErrorToCode),
		logging.WithFieldsFromContextAndCallMeta(logSroData),
	}

	return grpc.NewServer(
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
			grpc.ChainUnaryInterceptor(
				selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(sroauth.AuthFunc(kcClient, realm)), selector.MatchFunc(sroauth.NotPublicServiceMatcher)),
				logging.UnaryServerInterceptor(interceptorLogger(log.Logger), opts...),
			),
			grpc.ChainStreamInterceptor(
				selector.StreamServerInterceptor(auth.StreamServerInterceptor(sroauth.AuthFunc(kcClient, realm)), selector.MatchFunc(sroauth.NotPublicServiceMatcher)),
				logging.StreamServerInterceptor(interceptorLogger(log.Logger), opts...),
			)),
		runtime.NewServeMux()
}

// InterceptorLogger adapts logrus logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func interceptorLogger(l logrus.FieldLogger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make(map[string]any, len(fields)/2)
		i := logging.Fields(fields).Iterator()
		for i.Next() {
			k, v := i.At()
			f[k] = v
		}

		l := l.WithFields(f)

		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg)
		case logging.LevelInfo:
			l.Info(msg)
		case logging.LevelWarn:
			l.Warn(msg)
		case logging.LevelError:
			l.Error(msg)
		default:
			l.Debugf("unknown level %v: %s", lvl, msg)
		}
	})
}

func logSroData(ctx context.Context, c interceptors.CallMeta) logging.Fields {
	out := logging.Fields{}
	if spanCtx := trace.SpanContextFromContext(ctx); spanCtx.IsValid() {
		out = append(out, "traceId", spanCtx.TraceID().String())
	}
	if claims, ok := sroauth.RetrieveClaims(ctx); ok {
		out = append(out, "requestor", fmt.Sprintf("%s:%s", claims.Username, claims.Subject))
	}
	return out
}
