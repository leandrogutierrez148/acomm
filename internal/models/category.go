package models

type Category struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Code string `json:"code" gorm:"uniqueIndex;not null,default:null"`
	Name string `json:"name"`
}

func (c *Category) TableName() string {
	return "product_categories"
}
