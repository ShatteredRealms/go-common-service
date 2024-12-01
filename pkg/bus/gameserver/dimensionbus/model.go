package dimensionbus

import (
	"time"
)

type Dimension struct {
	Id        string `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Dimensions []*Dimension
