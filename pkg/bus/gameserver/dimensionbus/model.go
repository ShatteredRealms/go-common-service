package dimensionbus

import (
	"time"

	"github.com/google/uuid"
)

type Dimension struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Dimensions []*Dimension
