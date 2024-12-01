package mapbus

import "time"

type Map struct {
	Id        string `gorm:"primaryKey" json:"id"`
	OwnerId   string `gorm:"index;not null" json:"ownerId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Maps []*Map
