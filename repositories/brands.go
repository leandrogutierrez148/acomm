package repositories

import (
	"github.com/lgutierrez148/acomm/interfaces"
	"github.com/lgutierrez148/acomm/models"
)

type BrandsRepository struct {
	db interfaces.IDatabase
}

func NewBrandsRepository(db interfaces.IDatabase) *BrandsRepository {
	return &BrandsRepository{db: db}
}

func (r *BrandsRepository) Create(brand *models.Brand) error {
	return r.db.Create(brand).GetError()
}

func (r *BrandsRepository) FindByID(id uint) (*models.Brand, error) {
	var brand models.Brand
	err := r.db.First(&brand, id).GetError()
	return &brand, err
}

func (r *BrandsRepository) FindAll() ([]models.Brand, error) {
	var brands []models.Brand
	err := r.db.Find(&brands).GetError()
	return brands, err
}

func (r *BrandsRepository) Update(brand *models.Brand) error {
	return r.db.Save(brand).GetError()
}

func (r *BrandsRepository) Delete(id uint) error {
	return r.db.Delete(&models.Brand{}, id).GetError()
}
