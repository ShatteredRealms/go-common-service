package channelbus

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	DimensionId uuid.UUID `gorm:"index;not null" json:"dimensionId"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Channels []*Channel
