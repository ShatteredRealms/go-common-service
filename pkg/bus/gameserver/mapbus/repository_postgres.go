package mapbus

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/srospan"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type postgresRepository struct {
	gormdb *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	db.AutoMigrate(&Map{})
	return &postgresRepository{gormdb: db}
}

// Save implements MapRepository.
func (p *postgresRepository) Save(ctx context.Context, msg Message) error {
	m := Map{
		Id: msg.Id,
	}

	updateSpanWithMap(ctx, m.Id)
	return p.db(ctx).Save(&m).Error
}

// Delete implements MapRepository.
func (p *postgresRepository) Delete(ctx context.Context, mapId string) error {
	updateSpanWithMap(ctx, mapId)
	return p.db(ctx).Delete(&Map{}, "id = ?", mapId).Error
}

// GetById implements MapRepository.
func (p *postgresRepository) GetById(ctx context.Context, mapId string) (m *Map, _ error) {
	result := p.db(ctx).First(&m, "id = ?", mapId)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	updateSpanWithMap(ctx, mapId)
	return m, nil
}

// GetAll implements MapRepository.
func (p *postgresRepository) GetAll(ctx context.Context) (maps *Maps, _ error) {
	return maps, p.db(ctx).Find(maps).Error
}

func (p *postgresRepository) db(ctx context.Context) *gorm.DB {
	return p.gormdb.WithContext(ctx)
}

func updateSpanWithMap(ctx context.Context, mapId string) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.MapId(mapId),
	)
}
