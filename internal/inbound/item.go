package inbound

import "github.com/lgutierrez148/acomm/internal/models"

type CreateItemRequest struct {
	ProductID uint    `json:"product_id"`
	SKU       string  `json:"sku"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
}

func (r *CreateItemRequest) ToDomain() *models.Item {
	return &models.Item{
		ProductID: r.ProductID,
		SKU:       r.SKU,
		Price:     r.Price,
		Stock:     r.Stock,
	}
}

type UpdateItemRequest struct {
	ProductID uint    `json:"product_id"`
	SKU       string  `json:"sku"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
}

func (r *UpdateItemRequest) ToDomain() *models.Item {
	return &models.Item{
		ProductID: r.ProductID,
		SKU:       r.SKU,
		Price:     r.Price,
		Stock:     r.Stock,
	}
}
