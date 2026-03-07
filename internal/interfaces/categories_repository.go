package interfaces

import "github.com/lgutierrez148/acomm/internal/models"

type ICategoriesRepository interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(category *models.Category) error
}
