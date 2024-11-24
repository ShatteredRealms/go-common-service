package srv

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/pb"
	"github.com/ShatteredRealms/go-common-service/pkg/util"
	"github.com/WilSimpson/gocloak/v13"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type healthService struct {
	pb.UnimplementedHealthServiceServer
}

func NewHealthServiceServer() pb.HealthServiceServer {
	return &healthService{}
}

func (s *healthService) Health(context.Context, *emptypb.Empty) (*pb.HealthMessage, error) {
	return &pb.HealthMessage{Status: "ok"}, nil
}

func SetupHealthServer(ctx context.Context, cfg *config.BaseConfig) error {
	keycloakClient := gocloak.NewClient(cfg.Keycloak.BaseURL)
	grpcServer, gwmux := util.InitServerDefaults(keycloakClient, cfg.Keycloak.Realm)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	pb.RegisterHealthServiceServer(grpcServer, NewHealthServiceServer())
	err := pb.RegisterHealthServiceHandlerFromEndpoint(ctx, gwmux, cfg.Server.Address(), opts)
	if err != nil {
		return fmt.Errorf("register health service handler endpoint: %w", err)
	}
	return nil
}

