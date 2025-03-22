package channelbus

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/google/uuid"
)

type Repository interface {
	bus.BusMessageRepository[Message]

	GetById(ctx context.Context, channelId *uuid.UUID) (*Channel, error)

	GetAll(ctx context.Context) (*Channels, error)
	GetByDimensionId(ctx context.Context, ownerId *uuid.UUID) (*Channels, error)
}
