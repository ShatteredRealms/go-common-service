package dimensionbus

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
	db.AutoMigrate(&Dimension{})
	return &postgresRepository{gormdb: db}
}

// Save implements DimensionRepository.
func (p *postgresRepository) Save(ctx context.Context, msg Message) error {
	dimension := Dimension{
		Id: msg.Id,
	}

	updateSpanWithDimension(ctx, dimension.Id)
	return p.db(ctx).Save(&dimension).Error
}

// Delete implements DimensionRepository.
func (p *postgresRepository) Delete(ctx context.Context, dimensionId string) error {
	updateSpanWithDimension(ctx, dimensionId)
	return p.db(ctx).Delete(&Dimension{}, "id = ?", dimensionId).Error
}

// GetById implements DimensionRepository.
func (p *postgresRepository) GetById(ctx context.Context, dimensionId string) (dimension *Dimension, _ error) {
	result := p.db(ctx).First(&dimension, "id = ?", dimensionId)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	updateSpanWithDimension(ctx, dimensionId)
	return dimension, nil
}

// GetAll implements DimensionRepository.
func (p *postgresRepository) GetAll(ctx context.Context) (dimensions *Dimensions, _ error) {
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
