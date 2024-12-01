package mapbus

import (
	"context"
)

type Repository interface {
	GetMapById(ctx context.Context, dimensionId string) (*Map, error)

	GetMaps(ctx context.Context) (*Maps, error)

	CreateMap(ctx context.Context, dimensionId string) (*Map, error)

	DeleteMap(ctx context.Context, dimensionId string) (*Map, error)
}
