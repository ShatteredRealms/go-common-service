package mapbus

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
)

type Repository interface {
	bus.BusMessageRepository[Message]

	GetById(ctx context.Context, dimensionId string) (*Map, error)

	GetAll(ctx context.Context) (*Maps, error)
}
