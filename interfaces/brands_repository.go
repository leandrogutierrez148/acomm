package interfaces

import "github.com/lgutierrez148/acomm/models"

type IBrandsRepository interface {
	Create(brand *models.Brand) error
	FindByID(id uint) (*models.Brand, error)
	FindAll() ([]models.Brand, error)
	Update(brand *models.Brand) error
	Delete(id uint) error
}
