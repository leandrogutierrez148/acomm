package inbound

import "github.com/lgutierrez148/acomm/internal/models"

type CreateItemSpecificationRequest struct {
	ItemID uint   `json:"item_id"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

func (r *CreateItemSpecificationRequest) ToDomain() *models.ItemSpecification {
	return &models.ItemSpecification{
		ItemID: r.ItemID,
		Key:    r.Key,
		Value:  r.Value,
	}
}

type UpdateItemSpecificationRequest struct {
	ItemID uint   `json:"item_id"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

func (r *UpdateItemSpecificationRequest) ToDomain() *models.ItemSpecification {
	return &models.ItemSpecification{
		ItemID: r.ItemID,
		Key:    r.Key,
		Value:  r.Value,
	}
}
