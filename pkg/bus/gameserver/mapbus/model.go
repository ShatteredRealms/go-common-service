package mapbus

import (
	"time"

	"github.com/google/uuid"
)

type Map struct {
	Id        uuid.UUID `gorm:"primaryKey" json:"id"`
	OwnerId   uuid.UUID `gorm:"index;not null" json:"ownerId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Maps []*Map
