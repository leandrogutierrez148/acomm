package repositories

import (
	"github.com/lgutierrez148/acomm/models"
	"gorm.io/gorm"
)

type BrandsRepository struct {
	db *gorm.DB
}

func NewBrandsRepository(db *gorm.DB) *BrandsRepository {
	return &BrandsRepository{db: db}
}

func (r *BrandsRepository) Create(brand *models.Brand) error {
	return r.db.Create(brand).Error
}

func (r *BrandsRepository) FindByID(id uint) (*models.Brand, error) {
	var brand models.Brand
	err := r.db.First(&brand, id).Error
	return &brand, err
}

func (r *BrandsRepository) FindAll() ([]models.Brand, error) {
	var brands []models.Brand
	err := r.db.Find(&brands).Error
	return brands, err
}

func (r *BrandsRepository) Update(brand *models.Brand) error {
	return r.db.Save(brand).Error
}

func (r *BrandsRepository) Delete(id uint) error {
	return r.db.Delete(&models.Brand{}, id).Error
}
