package dimensionbus

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/google/uuid"
)

type Repository interface {
	bus.BusMessageRepository[Message]

	GetById(ctx context.Context, dimensionId *uuid.UUID) (*Dimension, error)

	GetAll(ctx context.Context) (*Dimensions, error)
}
