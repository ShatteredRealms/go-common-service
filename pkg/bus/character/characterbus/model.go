package characterbus

import (
	"time"

	"github.com/google/uuid"
)

type Character struct {
	Id          uuid.UUID `gorm:"primaryKey" json:"id"`
	OwnerId     uuid.UUID `gorm:"index;not null" json:"ownerId"`
	DimensionId uuid.UUID `gorm:"index;not null" json:"dimensionId"`
	MapId       uuid.UUID `gorm:"index;not null" json:"mapId"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Characters []*Character
