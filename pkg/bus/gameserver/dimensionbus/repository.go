package dimensionbus

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
)

type Repository interface {
	bus.BusMessageRepository

	GetById(ctx context.Context, dimensionId string) (*Dimension, error)

	GetAll(ctx context.Context) (*Dimensions, error)
}
