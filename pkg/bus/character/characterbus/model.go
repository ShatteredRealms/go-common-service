package characterbus

import (
	"time"

	"github.com/google/uuid"
)

type Character struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OwnerId     uuid.UUID `gorm:"index;not null;type:uuid" json:"ownerId"`
	DimensionId uuid.UUID `gorm:"index;not null;type:uuid" json:"dimensionId"`
	MapId       uuid.UUID `gorm:"index;not null;type:uuid" json:"mapId"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Characters []*Character

func (chars Characters) ToIds() []string {
	ids := make([]string, len(chars))
	for i, char := range chars {
		ids[i] = char.Id.String()
	}
	return ids
}
