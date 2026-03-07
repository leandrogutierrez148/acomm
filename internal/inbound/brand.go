package inbound

import "github.com/lgutierrez148/acomm/internal/models"

type CreateBrandRequest struct {
	Name               string `json:"name"`
	Title              string `json:"title"`
	ImageURL           string `json:"image_url"`
	MetaTagDescription string `json:"meta_tag_description"`
	IsActive           bool   `json:"is_active"`
}

func (r *CreateBrandRequest) ToDomain() *models.Brand {
	return &models.Brand{
		Name:               r.Name,
		Title:              r.Title,
		ImageURL:           r.ImageURL,
		MetaTagDescription: r.MetaTagDescription,
		IsActive:           r.IsActive,
	}
}

type UpdateBrandRequest struct {
	Name               string `json:"name"`
	Title              string `json:"title"`
	ImageURL           string `json:"image_url"`
	MetaTagDescription string `json:"meta_tag_description"`
	IsActive           bool   `json:"is_active"`
}

func (r *UpdateBrandRequest) ToDomain() *models.Brand {
	return &models.Brand{
		Name:               r.Name,
		Title:              r.Title,
		ImageURL:           r.ImageURL,
		MetaTagDescription: r.MetaTagDescription,
		IsActive:           r.IsActive,
	}
}
