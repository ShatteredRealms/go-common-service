package mapbus

import (
	"time"

	"github.com/google/uuid"
)

type Map struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OwnerId   uuid.UUID `gorm:"index;not null;type:uuid" json:"ownerId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Maps []*Map
