package repositories

import (
	"github.com/lgutierrez148/acomm/internal/interfaces"
	"github.com/lgutierrez148/acomm/internal/models"
)

type ProductsSpecificationsRepository struct {
	db interfaces.IDatabase
}

func NewProductsSpecificationsRepository(db interfaces.IDatabase) *ProductsSpecificationsRepository {
	return &ProductsSpecificationsRepository{db: db}
}

func (r *ProductsSpecificationsRepository) Create(spec *models.ProductSpecification) error {
	return r.db.Create(spec).GetError()
}

func (r *ProductsSpecificationsRepository) FindByID(id uint) (*models.ProductSpecification, error) {
	var spec models.ProductSpecification
	err := r.db.First(&spec, id).GetError()
	return &spec, err
}

func (r *ProductsSpecificationsRepository) FindAll() ([]models.ProductSpecification, error) {
	var specs []models.ProductSpecification
	err := r.db.Find(&specs).GetError()
	return specs, err
}

func (r *ProductsSpecificationsRepository) FindByProductID(productID uint) ([]models.ProductSpecification, error) {
	var specs []models.ProductSpecification
	err := r.db.Where("product_id = ?", productID).Find(&specs).GetError()
	return specs, err
}

func (r *ProductsSpecificationsRepository) Update(spec *models.ProductSpecification) error {
	return r.db.Save(spec).GetError()
}

func (r *ProductsSpecificationsRepository) Delete(id uint) error {
	return r.db.Delete(&models.ProductSpecification{}, id).GetError()
}
