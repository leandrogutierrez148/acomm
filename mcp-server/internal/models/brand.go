package models

// Brand represents a product brand.
type Brand struct {
	ID                 uint   `json:"id" gorm:"primaryKey"`
	ImageURL           string `json:"image_url"`
	IsActive           bool   `json:"is_active"`
	MetaTagDescription string `json:"meta_tag_description"`
	Name               string `json:"name" gorm:"uniqueIndex;not null"`
	Title              string `json:"title"`
}

func (b *Brand) TableName() string {
	return "brands"
}
