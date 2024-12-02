package characterbus

import "time"

type Character struct {
	Id          string `gorm:"primaryKey" json:"id"`
	OwnerId     string `gorm:"index;not null" json:"ownerId"`
	DimensionId string `gorm:"index;not null" json:"dimensionId"`
	MapId       string `gorm:"index;not null" json:"mapId"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Characters []*Character
