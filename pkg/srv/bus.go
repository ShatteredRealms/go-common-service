package srv

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/pb"
	"github.com/ShatteredRealms/go-common-service/pkg/util"
	"github.com/WilSimpson/gocloak/v13"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WriterResetCallback func() error

type busService struct {
	pb.UnimplementedBusServiceServer
	ctx Context

	readerResetters map[bus.BusMessageType]bus.Resettable
	writerCallbacks map[bus.BusMessageType]WriterResetCallback
}

var (
	BusRoles = []*gocloak.Role{}

	RoleBusReset = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("bus.reset"),
		Description: gocloak.StringP("Allow resetting the reader and writer buses"),
	}, &BusRoles)

	ErrBusNotFound = errors.New("bus not registered")
	ErrBusReset    = errors.New("bus reset failed")
)

// ResetReaderBus implements pb.BusServiceServer.
func (b *busService) ResetReaderBus(ctx context.Context, request *pb.BusTarget) (*pb.ResetBusResponse, error) {
	if request.GetType() == "" {
		// Reset all buses
		var err error
		builder := strings.Builder{}
		for name, reader := range b.readerResetters {
			err = errors.Join(err, reader.Reset(ctx))
			builder.WriteString(string(name))
			builder.WriteString(", ")
		}
		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Errorf("%w: %w", ErrBusReset, err).Error())
		}

		str := builder.String()
		return &pb.ResetBusResponse{
			Message: fmt.Sprintf("Reset %d buses: %s", len(b.readerResetters), str[:len(str)-2]),
		}, nil
	}

	// Reset a specific bus
	bus, ok := b.readerResetters[bus.BusMessageType(request.GetType())]
	if !ok {
		return nil, status.Errorf(codes.NotFound, ErrBusNotFound.Error())
	}

	err := bus.Reset(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Errorf("%w: %w", ErrBusReset, err).Error())
	}

	return &pb.ResetBusResponse{
		Message: fmt.Sprintf("Reset 1 bus: %s", request.GetType()),
	}, nil
}

// ResetWriterBus implements pb.BusServiceServer.
func (b *busService) ResetWriterBus(ctx context.Context, request *pb.BusTarget) (*pb.ResetBusResponse, error) {
	if request.GetType() == "" {
		var err error
		builder := strings.Builder{}
		for name, busCallback := range b.writerCallbacks {
			err = errors.Join(err, busCallback())
			builder.WriteString(string(name))
			builder.WriteString(", ")
		}

		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Errorf("%w: %w", ErrBusReset, err).Error())
		}
		str := builder.String()
		return &pb.ResetBusResponse{
			Message: fmt.Sprintf("Reset %d buses: %s", len(b.writerCallbacks), str[:len(str)-2]),
		}, nil
	}

	busCallback, ok := b.writerCallbacks[bus.BusMessageType(request.GetType())]
	if !ok {
		return nil, status.Errorf(codes.NotFound, ErrBusNotFound.Error())
	}

	err := busCallback()
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Errorf("%w: %w", ErrBusReset, err).Error())
	}

	return &pb.ResetBusResponse{
		Message: fmt.Sprintf("Reset 1 bus: %s", request.GetType()),
	}, nil
}

func NewBusServiceServer(
	ctx context.Context,
	srvCtx Context,
	readerResetters map[bus.BusMessageType]bus.Resettable,
	writerCallbacks map[bus.BusMessageType]WriterResetCallback,
) (*busService, error) {
	err := srvCtx.CreateRoles(ctx, &BusRoles)
	if err != nil {
		return nil, fmt.Errorf("failed to create roles: %w", err)
	}

	service := &busService{
		ctx:             srvCtx,
		writerCallbacks: writerCallbacks,
		readerResetters: readerResetters,
	}

	return service, nil
}
