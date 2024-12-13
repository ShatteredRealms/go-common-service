package dimensionbus

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/common"
	"github.com/google/uuid"
)

type Service interface {
	bus.BusProcessor[Message]
	GetDimensions(ctx context.Context) (*Dimensions, error)
	GetDimensionById(ctx context.Context, dimensionId string) (*Dimension, error)
}

type service struct {
	bus.DefaultBusProcessor[Message]
}

func NewService(
	repo Repository,
	dimensionBus bus.MessageBusReader[Message],
) Service {
	return &service{
		DefaultBusProcessor: bus.DefaultBusProcessor[Message]{
			Reader: dimensionBus,
			Repo:   repo,
		},
	}
}

// GetDimensionById implements DimensionService.
func (d *service) GetDimensionById(ctx context.Context, dimensionId string) (*Dimension, error) {
	id, err := uuid.Parse(dimensionId)
	if err != nil {
		return nil, common.ErrInvalidId
	}

	return d.Repo.(Repository).GetById(ctx, &id)
}

// GetDimensions implements DimensionService.
func (d *service) GetDimensions(ctx context.Context) (*Dimensions, error) {
	return d.Repo.(Repository).GetAll(ctx)
}
