package models

import "time"

type ItemSpecification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ItemID    uint      `json:"item_id" gorm:"not null;index"`
	Key       string    `json:"key" gorm:"not null"`
	Value     string    `json:"value" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ItemSpecification) TableName() string {
	return "item_specifications"
}
