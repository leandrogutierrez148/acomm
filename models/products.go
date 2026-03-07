package models

import (
	"time"
)

// Product represents a product in the catalog.
// It includes a unique code and a price.
type Product struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CategoryID uint      `json:"category_id" gorm:"index"`
	Category   Category  `gorm:"foreignKey:CategoryID"`
	BrandID    uint      `json:"brand_id,omitempty" gorm:"index"`
	Brand      Brand     `gorm:"foreignKey:BrandID"`
	Items      []Item    `json:"items,omitempty" gorm:"foreignKey:ProductID"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (p *Product) TableName() string {
	return "products"
}
