package items

import "time"

type GetItemsResponse struct {
	Items []Item `json:"items"`
}

type GetItemResponse struct {
	Item Item `json:"item"`
}

type CreateItemResponse struct {
	Item Item `json:"item"`
}

type UpdateItemResponse struct {
	Item Item `json:"item"`
}

type Item struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	SKU       string    `json:"sku"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
