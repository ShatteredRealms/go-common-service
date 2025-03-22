package channelbus

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	Id          uuid.UUID `gorm:"primaryKey" json:"id"`
	DimensionId uuid.UUID `gorm:"index;not null" json:"dimensionId"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Channels []*Channel
