package interfaces

import "github.com/lgutierrez148/acomm/internal/models"

type IItemsSpecificationsRepository interface {
	Create(spec *models.ItemSpecification) error
	FindByID(id uint) (*models.ItemSpecification, error)
	FindAll() ([]models.ItemSpecification, error)
	FindByItemID(itemID uint) ([]models.ItemSpecification, error)
	Update(spec *models.ItemSpecification) error
	Delete(id uint) error
}
