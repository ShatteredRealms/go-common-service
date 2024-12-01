package characterbus

import (
	"context"
)

type Repository interface {
	GetCharacterById(ctx context.Context, characterId string) (*Character, error)

	GetCharacters(ctx context.Context) (*Characters, error)
	GetCharactersByOwnerId(ctx context.Context, ownerId string) (*Characters, error)

	CreateCharacter(ctx context.Context, characterId, ownerId string) (*Character, error)

	UpdateCharacter(ctx context.Context, character *Character) (*Character, error)

	DeleteCharacter(ctx context.Context, characterId string) (*Character, error)

	DoesOwnCharacter(ctx context.Context, characterId, ownerId string) (bool, error)
}
