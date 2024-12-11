package mapbus

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
)

type Service[T bus.BusMessage[any]] interface {
	bus.BusProcessor[T]
	GetMaps(ctx context.Context) (*Maps, error)
	GetMapById(ctx context.Context, mapId string) (*Map, error)
}

type service struct {
	bus.DefaultBusProcessor[Message]
}

func NewService(
	repo Repository,
	mapBus bus.MessageBusReader[Message],
) Service[Message] {
	return &service{
		DefaultBusProcessor: bus.DefaultBusProcessor[Message]{
			Reader: mapBus,
			Repo:   repo,
		},
	}
}

// GetMapById implements MapService.
func (d *service) GetMapById(ctx context.Context, mapId string) (*Map, error) {
	return d.Repo.(Repository).GetById(ctx, mapId)
}

// GetMaps implements MapService.
func (d *service) GetMaps(ctx context.Context) (*Maps, error) {
	return d.Repo.(Repository).GetAll(ctx)
}
