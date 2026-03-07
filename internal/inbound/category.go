package inbound

import "github.com/lgutierrez148/acomm/internal/models"

type CreateCategoryRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (r *CreateCategoryRequest) ToDomain() *models.Category {
	return &models.Category{
		Name: r.Name,
		Code: r.Code,
	}
}
