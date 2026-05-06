package models

import "time"

type ProductSpecification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id" gorm:"not null;index"`
	Key       string    `json:"key" gorm:"not null"`
	Value     string    `json:"value" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ProductSpecification) TableName() string {
	return "product_specifications"
}
