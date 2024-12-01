package dimensionbus

import (
	"context"
)

type Repository interface {
	GetDimensionById(ctx context.Context, dimensionId string) (*Dimension, error)

	GetDimensions(ctx context.Context) (*Dimensions, error)

	CreateDimension(ctx context.Context, dimensionId string) (*Dimension, error)

	DeleteDimension(ctx context.Context, dimensionId string) (*Dimension, error)
}
