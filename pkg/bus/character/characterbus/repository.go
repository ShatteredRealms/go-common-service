package characterbus

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/google/uuid"
)

type Repository interface {
	bus.BusMessageRepository[Message]

	GetById(ctx context.Context, characterId *uuid.UUID) (*Character, error)

	GetAll(ctx context.Context) (*Characters, error)
	GetByOwnerId(ctx context.Context, ownerId *uuid.UUID) (*Characters, error)

	IsOwner(ctx context.Context, characterId, ownerId *uuid.UUID) (bool, error)
}
