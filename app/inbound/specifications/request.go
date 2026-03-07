package specifications

import "github.com/lgutierrez148/acomm/models"

type CreateSpecificationRequest struct {
	ProductID uint   `json:"product_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

func (r *CreateSpecificationRequest) ToDomain() *models.Specification {
	return &models.Specification{
		ProductID: r.ProductID,
		Key:       r.Key,
		Value:     r.Value,
	}
}

type UpdateSpecificationRequest struct {
	ProductID uint   `json:"product_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

func (r *UpdateSpecificationRequest) ToDomain() *models.Specification {
	return &models.Specification{
		ProductID: r.ProductID,
		Key:       r.Key,
		Value:     r.Value,
	}
}
