package models

import "time"

// Order represents a purchase containing multiple items.
type Order struct {
	ID              int         `json:"id" gorm:"primaryKey;column:id"`
	Items           []ItemOrder `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	CustomerEmail   string      `json:"customer_email,omitempty"`
	CustomerName    string      `json:"customer_name,omitempty"`
	CustomerAddress string      `json:"customer_address,omitempty"`
	CustomerPhone   string      `json:"customer_phone,omitempty"`
	Status          string      `json:"status,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

func (o *Order) TableName() string {
	return "orders"
}

// ItemOrder represents an item within an Order.
type ItemOrder struct {
	ItemID  int `json:"item_id" gorm:"primaryKey"`
	OrderID int `json:"order_id" gorm:"primaryKey;not null"`
}

func (io *ItemOrder) TableName() string {
	return "item_orders"
}
