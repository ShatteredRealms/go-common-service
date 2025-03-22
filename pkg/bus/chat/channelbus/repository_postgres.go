package channelbus

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
	ErrNilChannel = errors.New("channel is nil")
)

type postgresRepository struct {
	gormdb *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	db.AutoMigrate(&Channel{})
	return &postgresRepository{gormdb: db}
}

// SaveChannel implements ChannelRepository.
func (p *postgresRepository) Save(
	ctx context.Context,
	msg Message,
) error {
	channel := Channel{
		Id:          msg.Id,
		DimensionId: msg.DimensionId,
	}

	updateSpanWithChannel(ctx, channel.Id.String())
	return p.db(ctx).Save(&channel).Error
}

// DeleteChannel implements ChannelRepository.
func (p *postgresRepository) Delete(
	ctx context.Context,
	id *uuid.UUID,
) error {
	channel := &Channel{}
	err := p.db(ctx).Clauses(clause.Returning{}).Delete(channel, "id = ?", id).Error
	if err != nil {
		return err
	}

	updateSpanWithChannel(ctx, id.String())
	return err
}

// GetById implements ChannelRepository.
func (p *postgresRepository) GetById(
	ctx context.Context,
	channelId *uuid.UUID,
) (channel *Channel, _ error) {
	channel = &Channel{}
	result := p.db(ctx).First(&channel, "id = ?", channelId)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	updateSpanWithChannel(ctx, channelId.String())
	return channel, nil
}

// GetAll implements ChannelRepository.
func (p *postgresRepository) GetAll(
	ctx context.Context,
) (channels *Channels, _ error) {
	return channels, p.db(ctx).Find(&channels).Error
}

// GetByDimensionId implements ChannelRepository.
func (p *postgresRepository) GetByDimensionId(
	ctx context.Context,
	ownerId *uuid.UUID,
) (channels *Channels, _ error) {
	return channels, p.db(ctx).Where("dimension_id = ?", ownerId).Find(&channels).Error
}

// IsOwner implements ChannelRepository.
func (p *postgresRepository) IsOwner(ctx context.Context, channelId *uuid.UUID, ownerId *uuid.UUID) (bool, error) {
	result := p.db(ctx).Where("id = ? AND owner_id = ?", channelId, ownerId).First(&Channel{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (p *postgresRepository) db(ctx context.Context) *gorm.DB {
	return p.gormdb.WithContext(ctx)
}

func updateSpanWithChannel(ctx context.Context, channelId string) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.ChatChannelId(channelId),
	)
}
