package outbound

import "time"

type GetProductSpecificationsResponse struct {
	Specifications []ProductSpecification `json:"specifications"`
}

type GetProductSpecificationResponse struct {
	Specification ProductSpecification `json:"specification"`
}

type CreateProductSpecificationResponse struct {
	Specification ProductSpecification `json:"specification"`
}

type UpdateProductSpecificationResponse struct {
	Specification ProductSpecification `json:"specification"`
}

type ProductSpecification struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
