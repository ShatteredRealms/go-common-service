package characterbus

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
)

type Service interface {
	bus.BusProcessor
	GetCharacters(ctx context.Context) (*Characters, error)
	GetCharacterById(ctx context.Context, characterId string) (*Character, error)
	DoesOwnCharacter(ctx context.Context, characterId, ownerId string) (bool, error)
}

type service struct {
	bus.DefaultBusProcessor[Message]
}

func NewService(
	repo Repository,
	characterBus bus.MessageBusReader[Message],
) Service {
	return &service{
		DefaultBusProcessor: bus.DefaultBusProcessor[Message]{
			Reader: characterBus,
			Repo:   repo,
		},
	}
}

// GetCharacterById implements CharacterService.
func (d *service) GetCharacterById(ctx context.Context, characterId string) (*Character, error) {
	return d.Repo.(Repository).GetById(ctx, characterId)
}

// GetCharacters implements CharacterService.
func (d *service) GetCharacters(ctx context.Context) (*Characters, error) {
	return d.Repo.(Repository).GetAll(ctx)
}

// DoesOwnCharacter implements Service.
func (d *service) DoesOwnCharacter(ctx context.Context, characterId string, ownerId string) (bool, error) {
	return d.Repo.(Repository).IsOwner(ctx, characterId, ownerId)
}
