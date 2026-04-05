package interfaces

import "github.com/leandrogutierrez148/acomm/internal/models"

type ICategoriesRepository interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(category *models.Category) error
}
