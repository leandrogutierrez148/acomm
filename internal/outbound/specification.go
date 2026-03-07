package outbound

import "time"

type GetSpecificationsResponse struct {
	Specifications []Specification `json:"specifications"`
}

type GetSpecificationResponse struct {
	Specification Specification `json:"specification"`
}

type CreateSpecificationResponse struct {
	Specification Specification `json:"specification"`
}

type UpdateSpecificationResponse struct {
	Specification Specification `json:"specification"`
}

type Specification struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
