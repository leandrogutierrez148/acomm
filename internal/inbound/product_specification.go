package inbound

import "github.com/lgutierrez148/acomm/internal/models"

type CreateProductSpecificationRequest struct {
	ProductID uint   `json:"product_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

func (r *CreateProductSpecificationRequest) ToDomain() *models.ProductSpecification {
	return &models.ProductSpecification{
		ProductID: r.ProductID,
		Key:       r.Key,
		Value:     r.Value,
	}
}

type UpdateProductSpecificationRequest struct {
	ProductID uint   `json:"product_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

func (r *UpdateProductSpecificationRequest) ToDomain() *models.ProductSpecification {
	return &models.ProductSpecification{
		ProductID: r.ProductID,
		Key:       r.Key,
		Value:     r.Value,
	}
}
