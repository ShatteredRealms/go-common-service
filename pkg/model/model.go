package model

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	Id        uuid.UUID  `db:"id" json:"id"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt" mapstructure:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt" mapstructure:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt" mapstructure:"deleted_at"`
}
