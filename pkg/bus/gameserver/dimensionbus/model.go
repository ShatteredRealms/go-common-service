package dimensionbus

import (
	"time"

	"github.com/google/uuid"
)

type Dimension struct {
	Id        uuid.UUID `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Dimensions []*Dimension
