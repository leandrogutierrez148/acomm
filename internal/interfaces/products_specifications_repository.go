package interfaces

import "github.com/lgutierrez148/acomm/internal/models"

type IProductsSpecificationsRepository interface {
	Create(spec *models.ProductSpecification) error
	FindByID(id uint) (*models.ProductSpecification, error)
	FindAll() ([]models.ProductSpecification, error)
	FindByProductID(productID uint) ([]models.ProductSpecification, error)
	Update(spec *models.ProductSpecification) error
	Delete(id uint) error
}
