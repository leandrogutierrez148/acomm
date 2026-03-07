package repositories

import (
	"github.com/lgutierrez148/acomm/models"
	"gorm.io/gorm"
)

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{
		db: db,
	}
}

func (r *CategoriesRepository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoriesRepository) CreateCategory(category *models.Category) error {
	if err := r.db.Create(category).Error; err != nil {
		return err
	}
	return nil
}
