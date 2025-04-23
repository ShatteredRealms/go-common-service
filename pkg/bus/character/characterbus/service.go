package characterbus

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/common"
	"github.com/google/uuid"
)

type Service interface {
	bus.BusProcessor[Message]
	GetCharacters(ctx context.Context) (Characters, error)
	GetCharacterById(ctx context.Context, characterId string) (*Character, error)
	GetCharactersByOwnerId(ctx context.Context, ownerId string) (Characters, error)
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
	id, err := uuid.Parse(characterId)
	if err != nil {
		return nil, common.ErrInvalidId
	}

	return d.Repo.(Repository).GetById(ctx, &id)
}

// GetCharacterById implements CharacterService.
func (d *service) GetCharactersByOwnerId(ctx context.Context, ownerId string) (Characters, error) {
	id, err := uuid.Parse(ownerId)
	if err != nil {
		return nil, common.ErrInvalidId
	}

	return d.Repo.(Repository).GetByOwnerId(ctx, &id)
}

// GetCharacters implements CharacterService.
func (d *service) GetCharacters(ctx context.Context) (Characters, error) {
	return d.Repo.(Repository).GetAll(ctx)
}

// DoesOwnCharacter implements Service.
func (d *service) DoesOwnCharacter(ctx context.Context, characterId string, ownerId string) (bool, error) {
	cId, err := uuid.Parse(characterId)
	if err != nil {
		return false, common.ErrInvalidId
	}
	oId, err := uuid.Parse(ownerId)
	if err != nil {
		return false, common.ErrInvalidId
	}

	return d.Repo.(Repository).IsOwner(ctx, &cId, &oId)
}
