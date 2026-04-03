package interfaces

import "github.com/lgutierrez148/acomm/internal/models"

type IItemImagesRepository interface {
	Create(image *models.ItemImage) error
	FindByID(id uint) (*models.ItemImage, error)
	FindAll() ([]models.ItemImage, error)
	FindByItemID(itemID uint) ([]models.ItemImage, error)
	Update(image *models.ItemImage) error
	Delete(id uint) error
}
