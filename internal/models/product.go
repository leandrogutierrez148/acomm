package models

import (
	"time"
)

// Product represents a product in the catalog.
type Product struct {
	ID                     uint      `json:"id" gorm:"primaryKey"`
	Name                   string    `json:"name"`
	DepartmentID           uint      `json:"department_id"`
	CategoryID             uint      `json:"category_id" gorm:"index"`
	Category               Category  `gorm:"foreignKey:CategoryID"`
	BrandID                uint      `json:"brand_id,omitempty" gorm:"index"`
	Brand                  Brand     `gorm:"foreignKey:BrandID"`
	LinkID                 string    `json:"link_id"`
	RefID                  string    `json:"ref_id"`
	IsVisible              bool      `json:"is_visible"`
	Description            string    `json:"description"`
	DescriptionShort       string    `json:"description_short"`
	ReleaseDate            time.Time `json:"release_date"`
	Keywords               string    `json:"keywords"`
	Title                  string    `json:"title"`
	IsActive               bool      `json:"is_active"`
	TaxCode                string    `json:"tax_code"`
	MetaTagDescription     string    `json:"meta_tag_description"`
	SupplierID             *uint     `json:"supplier_id"`
	ShowWithoutStock       bool      `json:"show_without_stock"`
	AdWordsRemarketingCode string    `json:"ad_words_remarketing_code"`
	LomadeeCampaignCode    string    `json:"lomadee_campaign_code"`
	Score                  int       `json:"score"`
	Items                  []Item    `json:"items,omitempty" gorm:"foreignKey:ProductID"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

func (p *Product) TableName() string {
	return "products"
}
