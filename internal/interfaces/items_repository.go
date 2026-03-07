package interfaces

import "github.com/lgutierrez148/acomm/internal/models"

type IItemsRepository interface {
	Create(item *models.Item) error
	FindByID(id uint) (*models.Item, error)
	FindAll() ([]models.Item, error)
	FindByProductID(productID uint) ([]models.Item, error)
	Update(item *models.Item) error
	Delete(id uint) error
}
