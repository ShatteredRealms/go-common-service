package util

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func GRPCHandlerFunc(grpcServer http.Handler, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType, x-grpc-web")

			if r.Method == "OPTIONS" {
				return
			}
			otelhttp.WithRouteTag(r.URL.Path, otherHandler).ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func StartServer(
	ctx context.Context,
	grpcServer http.Handler,
	gwmux http.Handler,
	address string,
) (*http.Server, chan error) {
	log.Logger.WithContext(ctx).Info("Starting server")
	listen, err := net.Listen("tcp", address)
	ch := make(chan error, 1)
	if err != nil {
		go func() {
			ch <- fmt.Errorf("listen server: %w", err)
		}()
		return nil, ch
	}

	httpSrv := &http.Server{
		Addr:    address,
		Handler: GRPCHandlerFunc(grpcServer, gwmux),
	}
	go func() {
		ch <- httpSrv.Serve(listen)
	}()
	return httpSrv, ch
}
