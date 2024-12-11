package mapbus

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
)

type Service interface {
	bus.BusProcessor[Message]
	GetMaps(ctx context.Context) (*Maps, error)
	GetMapById(ctx context.Context, mapId string) (*Map, error)
}

type service struct {
	bus.DefaultBusProcessor[Message]
}

func NewService(
	repo Repository,
	mapBus bus.MessageBusReader[Message],
) Service {
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
