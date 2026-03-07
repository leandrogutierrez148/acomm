package repositories

import (
	"github.com/lgutierrez148/acomm/interfaces"
	"github.com/lgutierrez148/acomm/models"
)

type CategoriesRepository struct {
	db interfaces.IDatabase
}

func NewCategoriesRepository(db interfaces.IDatabase) *CategoriesRepository {
	return &CategoriesRepository{
		db: db,
	}
}

func (r *CategoriesRepository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Find(&categories).GetError(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoriesRepository) CreateCategory(category *models.Category) error {
	if err := r.db.Create(category).GetError(); err != nil {
		return err
	}
	return nil
}
