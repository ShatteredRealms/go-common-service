package mapbus

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/srospan"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type postgresRepository struct {
	gormdb *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	db.AutoMigrate(&Map{})
	return &postgresRepository{gormdb: db}
}

// CreateMap implements MapRepository.
func (p *postgresRepository) CreateMap(ctx context.Context, mId string) (m *Map, _ error) {
	updateSpanWithMap(ctx, mId)
	m.Id = mId
	return m, p.db(ctx).Create(m).Error
}

// DeleteMap implements MapRepository.
func (p *postgresRepository) DeleteMap(ctx context.Context, mId string) (m *Map, _ error) {
	updateSpanWithMap(ctx, mId)
	return m, p.db(ctx).Clauses(clause.Returning{}).Delete(m, "id = ?", mId).Error
}

// GetMapById implements MapRepository.
func (p *postgresRepository) GetMapById(ctx context.Context, mId string) (m *Map, _ error) {
	result := p.db(ctx).Where("id = ?", mId).Find(&m)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	updateSpanWithMap(ctx, mId)
	return m, nil
}

// GetMaps implements MapRepository.
func (p *postgresRepository) GetMaps(ctx context.Context) (maps *Maps, _ error) {
	return maps, p.db(ctx).Find(maps).Error
}

func (p *postgresRepository) db(ctx context.Context) *gorm.DB {
	return p.gormdb.WithContext(ctx)
}

func updateSpanWithMap(ctx context.Context, mId string) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.MapId(mId),
	)
}
