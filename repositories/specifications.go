package repositories

import (
	"github.com/lgutierrez148/acomm/interfaces"
	"github.com/lgutierrez148/acomm/models"
)

type SpecificationsRepository struct {
	db interfaces.IDatabase
}

func NewSpecificationsRepository(db interfaces.IDatabase) *SpecificationsRepository {
	return &SpecificationsRepository{db: db}
}

func (r *SpecificationsRepository) Create(spec *models.Specification) error {
	return r.db.Create(spec).GetError()
}

func (r *SpecificationsRepository) FindByID(id uint) (*models.Specification, error) {
	var spec models.Specification
	err := r.db.First(&spec, id).GetError()
	return &spec, err
}

func (r *SpecificationsRepository) FindAll() ([]models.Specification, error) {
	var specs []models.Specification
	err := r.db.Find(&specs).GetError()
	return specs, err
}

func (r *SpecificationsRepository) FindByProductID(productID uint) ([]models.Specification, error) {
	var specs []models.Specification
	err := r.db.Where("product_id = ?", productID).Find(&specs).GetError()
	return specs, err
}

func (r *SpecificationsRepository) Update(spec *models.Specification) error {
	return r.db.Save(spec).GetError()
}

func (r *SpecificationsRepository) Delete(id uint) error {
	return r.db.Delete(&models.Specification{}, id).GetError()
}
