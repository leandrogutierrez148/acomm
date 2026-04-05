package inbound

import "github.com/leandrogutierrez148/acomm/internal/models"

type CreateItemImageRequest struct {
	Url    string  `json:"url"`
	Label  *string `json:"label"`
	Text   *string `json:"text"`
	IsMain bool    `json:"is_main"`
}

func (r *CreateItemImageRequest) ToDomain() *models.ItemImage {
	return &models.ItemImage{
		Url:    r.Url,
		Label:  r.Label,
		Text:   r.Text,
		IsMain: r.IsMain,
	}
}

type UpdateItemImageRequest struct {
	Url    string  `json:"url"`
	Label  *string `json:"label"`
	Text   *string `json:"text"`
	IsMain bool    `json:"is_main"`
}

func (r *UpdateItemImageRequest) ToDomain() *models.ItemImage {
	return &models.ItemImage{
		Url:    r.Url,
		Label:  r.Label,
		Text:   r.Text,
		IsMain: r.IsMain,
	}
}
