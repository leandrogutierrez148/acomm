package interfaces

import "github.com/lgutierrez148/acomm/models"

type ISpecificationsRepository interface {
	Create(spec *models.Specification) error
	FindByID(id uint) (*models.Specification, error)
	FindAll() ([]models.Specification, error)
	FindByProductID(productID uint) ([]models.Specification, error)
	Update(spec *models.Specification) error
	Delete(id uint) error
}
