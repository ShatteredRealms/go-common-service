package characterbus

import (
	"context"
	"errors"

	"github.com/ShatteredRealms/go-common-service/pkg/srospan"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrNilCharacter = errors.New("character is nil")
)

type postgresRepository struct {
	gormdb *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	db.AutoMigrate(&Character{})
	return &postgresRepository{gormdb: db}
}

// CreateCharacter implements CharacterRepository.
func (p *postgresRepository) CreateCharacter(
	ctx context.Context,
	characterId string,
	ownerId string,
) (*Character, error) {
	updateSpanWithCharacter(ctx, characterId, ownerId)
	character := &Character{
		Id:      characterId,
		OwnerId: ownerId,
	}
	return character, p.db(ctx).Create(&character).Error
}

// UpdateCharacter implements CharacterRepository.
func (p *postgresRepository) UpdateCharacter(
	ctx context.Context,
	character *Character,
) (*Character, error) {
	if character == nil {
		return nil, ErrNilCharacter
	}

	updateSpanWithCharacter(ctx, character.Id, character.OwnerId)
	return character, p.db(ctx).Save(&character).Error
}

// DeleteCharacter implements CharacterRepository.
func (p *postgresRepository) DeleteCharacter(
	ctx context.Context,
	characterId string,
) (character *Character, _ error) {
	err := p.db(ctx).Clauses(clause.Returning{}).Delete(&character, characterId).Error
	if err != nil {
		return nil, err
	}

	updateSpanWithCharacter(ctx, characterId, character.OwnerId)
	return character, err
}

// GetCharacterById implements CharacterRepository.
func (p *postgresRepository) GetCharacterById(
	ctx context.Context,
	characterId string,
) (character *Character, _ error) {
	result := p.db(ctx).Find(&character, characterId)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	updateSpanWithCharacter(ctx, characterId, character.OwnerId)
	return character, nil
}

// GetCharacters implements CharacterRepository.
func (p *postgresRepository) GetCharacters(
	ctx context.Context,
) (characters *Characters, _ error) {
	return characters, p.db(ctx).Find(&characters).Error
}

// GetCharacterByOwnerId implements CharacterRepository.
func (p *postgresRepository) GetCharactersByOwnerId(
	ctx context.Context,
	ownerId string,
) (characters *Characters, _ error) {
	return characters, p.db(ctx).Where("owner_id = ?", ownerId).Find(&characters).Error
}

// DoesOwnCharacter implements CharacterRepository.
func (p *postgresRepository) DoesOwnCharacter(ctx context.Context, characterId string, ownerId string) (bool, error) {
	result := p.db(ctx).Where("id = ? AND owner_id = ?", characterId, ownerId).Find(&Character{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (p *postgresRepository) db(ctx context.Context) *gorm.DB {
	return p.gormdb.WithContext(ctx)
}

func updateSpanWithCharacter(ctx context.Context, characterId string, ownerId string) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.TargetCharacterId(characterId),
		srospan.TargetOwnerId(ownerId),
	)
}
