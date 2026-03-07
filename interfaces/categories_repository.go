package interfaces

import "github.com/lgutierrez148/acomm/models"

type ICategoriesRepository interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(category *models.Category) error
}
