package channelbus

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/common"
	"github.com/google/uuid"
)

type Service interface {
	bus.BusProcessor[Message]
	GetChannels(ctx context.Context) (*Channels, error)
	GetChannelById(ctx context.Context, channelId string) (*Channel, error)
}

type service struct {
	bus.DefaultBusProcessor[Message]
}

func NewService(
	repo Repository,
	channelBus bus.MessageBusReader[Message],
) Service {
	return &service{
		DefaultBusProcessor: bus.DefaultBusProcessor[Message]{
			Reader: channelBus,
			Repo:   repo,
		},
	}
}

// GetChannelById implements ChannelService.
func (d *service) GetChannelById(ctx context.Context, channelId string) (*Channel, error) {
	id, err := uuid.Parse(channelId)
	if err != nil {
		return nil, common.ErrInvalidId
	}

	return d.Repo.(Repository).GetById(ctx, &id)
}

// GetChannels implements ChannelService.
func (d *service) GetChannels(ctx context.Context) (*Channels, error) {
	return d.Repo.(Repository).GetAll(ctx)
}
