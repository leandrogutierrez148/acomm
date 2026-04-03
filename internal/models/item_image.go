package models

import "time"

type ItemImage struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ItemID       uint      `json:"item_id" gorm:"not null;index"`
	ArchiveId    int       `json:"archive_id"`
	SkuId        int       `json:"sku_id"`
	Name         string    `json:"name"`
	Url          string    `json:"url" gorm:"not null"`
	FileLocation string    `json:"file_location"`
	Label        *string   `json:"label"`
	Text         *string   `json:"text"`
	IsMain       bool      `json:"is_main" gorm:"not null;default:false"`
	Position     int       `json:"position"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (ItemImage) TableName() string {
	return "item_images"
}
