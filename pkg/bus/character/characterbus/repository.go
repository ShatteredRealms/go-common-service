package characterbus

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
)

type Repository interface {
	bus.BusMessageRepository

	GetById(ctx context.Context, characterId string) (*Character, error)

	GetAll(ctx context.Context) (*Characters, error)
	GetByOwnerId(ctx context.Context, ownerId string) (*Characters, error)

	IsOwner(ctx context.Context, characterId, ownerId string) (bool, error)
}
