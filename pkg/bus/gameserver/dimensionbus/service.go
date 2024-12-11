package dimensionbus

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
)

type Service[T bus.BusMessage[any]] interface {
	bus.BusProcessor[T]
	GetDimensions(ctx context.Context) (*Dimensions, error)
	GetDimensionById(ctx context.Context, dimensionId string) (*Dimension, error)
}

type service struct {
	bus.DefaultBusProcessor[Message]
}

func NewService(
	repo Repository,
	dimensionBus bus.MessageBusReader[Message],
) Service[Message] {
	return &service{
		DefaultBusProcessor: bus.DefaultBusProcessor[Message]{
			Reader: dimensionBus,
			Repo:   repo,
		},
	}
}

// GetDimensionById implements DimensionService.
func (d *service) GetDimensionById(ctx context.Context, dimensionId string) (*Dimension, error) {
	return d.Repo.(Repository).GetById(ctx, dimensionId)
}

// GetDimensions implements DimensionService.
func (d *service) GetDimensions(ctx context.Context) (*Dimensions, error) {
	return d.Repo.(Repository).GetAll(ctx)
}
