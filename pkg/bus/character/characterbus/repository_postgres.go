package characterbus

import (
	"context"
	"errors"

	"github.com/ShatteredRealms/go-common-service/pkg/srospan"
	"github.com/google/uuid"
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

// SaveCharacter implements CharacterRepository.
func (p *postgresRepository) Save(
	ctx context.Context,
	msg Message,
) error {
	character := Character{
		Id:          msg.Id,
		OwnerId:     msg.OwnerId,
		DimensionId: msg.DimensionId,
		MapId:       msg.MapId,
	}

	updateSpanWithCharacter(ctx, character.Id.String(), character.OwnerId.String())
	return p.db(ctx).Save(&character).Error
}

// DeleteCharacter implements CharacterRepository.
func (p *postgresRepository) Delete(
	ctx context.Context,
	id *uuid.UUID,
) error {
	character := &Character{}
	err := p.db(ctx).Clauses(clause.Returning{}).Delete(character, "id = ?", id).Error
	if err != nil {
		return err
	}

	updateSpanWithCharacter(ctx, id.String(), character.OwnerId.String())
	return err
}

// GetById implements CharacterRepository.
func (p *postgresRepository) GetById(
	ctx context.Context,
	characterId *uuid.UUID,
) (character *Character, _ error) {
	character = &Character{}
	result := p.db(ctx).First(&character, "id = ?", characterId)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	updateSpanWithCharacter(ctx, characterId.String(), character.OwnerId.String())
	return character, nil
}

// GetAll implements CharacterRepository.
func (p *postgresRepository) GetAll(
	ctx context.Context,
) (characters Characters, _ error) {
	return characters, p.db(ctx).Find(&characters).Error
}

// GetCharacterByOwnerId implements CharacterRepository.
func (p *postgresRepository) GetByOwnerId(
	ctx context.Context,
	ownerId *uuid.UUID,
) (characters Characters, _ error) {
	return characters, p.db(ctx).Where("owner_id = ?", ownerId).Find(&characters).Error
}

// IsOwner implements CharacterRepository.
func (p *postgresRepository) IsOwner(ctx context.Context, characterId *uuid.UUID, ownerId *uuid.UUID) (bool, error) {
	result := p.db(ctx).Where("id = ? AND owner_id = ?", characterId, ownerId).First(&Character{})
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
