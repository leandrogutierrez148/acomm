package outbound

import "time"

type GetItemSpecificationsResponse struct {
	Specifications []ItemSpecification `json:"specifications"`
}

type GetItemSpecificationResponse struct {
	Specification ItemSpecification `json:"specification"`
}

type CreateItemSpecificationResponse struct {
	Specification ItemSpecification `json:"specification"`
}

type UpdateItemSpecificationResponse struct {
	Specification ItemSpecification `json:"specification"`
}

type ItemSpecification struct {
	ID        uint      `json:"id"`
	ItemID    uint      `json:"item_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
