package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/plugin/soft_delete"
)

type Model struct {
	Id        *uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:idx_deleted"`
}

func (m *Model) IsCreated() bool {
	return m.Id != nil
}
