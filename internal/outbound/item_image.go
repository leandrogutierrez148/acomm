package outbound

import "time"

type GetItemImagesResponse struct {
	Images []ItemImage `json:"images"`
}

type GetItemImageResponse struct {
	Image ItemImage `json:"image"`
}

type CreateItemImageResponse struct {
	Image ItemImage `json:"image"`
}

type UpdateItemImageResponse struct {
	Image ItemImage `json:"image"`
}

type ItemImage struct {
	ID        uint      `json:"id"`
	ItemID    uint      `json:"item_id"`
	Url       string    `json:"url"`
	Label     *string   `json:"label,omitempty"`
	Text      *string   `json:"text,omitempty"`
	IsMain    bool      `json:"is_main"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
