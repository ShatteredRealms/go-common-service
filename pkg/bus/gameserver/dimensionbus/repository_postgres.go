package dimensionbus

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
	db.AutoMigrate(&Dimension{})
	return &postgresRepository{gormdb: db}
}

// CreateDimension implements DimensionRepository.
func (p *postgresRepository) CreateDimension(ctx context.Context, dimensionId string) (dimension *Dimension, _ error) {
	updateSpanWithDimension(ctx, dimensionId)
	dimension.Id = dimensionId
	return dimension, p.db(ctx).Create(dimension).Error
}

// DeleteDimension implements DimensionRepository.
func (p *postgresRepository) DeleteDimension(ctx context.Context, dimensionId string) (dimension *Dimension, _ error) {
	updateSpanWithDimension(ctx, dimensionId)
	return dimension, p.db(ctx).Clauses(clause.Returning{}).Delete(dimension, "id = ?", dimensionId).Error
}

// GetDimensionById implements DimensionRepository.
func (p *postgresRepository) GetDimensionById(ctx context.Context, dimensionId string) (dimension *Dimension, _ error) {
	result := p.db(ctx).Where("id = ?", dimensionId).Find(&dimension)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	updateSpanWithDimension(ctx, dimensionId)
	return dimension, nil
}

// GetDimensions implements DimensionRepository.
func (p *postgresRepository) GetDimensions(ctx context.Context) (dimensions *Dimensions, _ error) {
	return dimensions, p.db(ctx).Find(dimensions).Error
}

func (p *postgresRepository) db(ctx context.Context) *gorm.DB {
	return p.gormdb.WithContext(ctx)
}

func updateSpanWithDimension(ctx context.Context, dimensionId string) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.DimensionId(dimensionId),
	)
}
