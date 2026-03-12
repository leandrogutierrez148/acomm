package repositories

import (
	"github.com/lgutierrez148/acomm/internal/interfaces"
	"github.com/lgutierrez148/acomm/internal/models"
)

type ItemsSpecificationsRepository struct {
	db interfaces.IDatabase
}

func NewItemsSpecificationsRepository(db interfaces.IDatabase) *ItemsSpecificationsRepository {
	return &ItemsSpecificationsRepository{db: db}
}

func (r *ItemsSpecificationsRepository) Create(spec *models.ItemSpecification) error {
	return r.db.Create(spec).GetError()
}

func (r *ItemsSpecificationsRepository) FindByID(id uint) (*models.ItemSpecification, error) {
	var spec models.ItemSpecification
	err := r.db.First(&spec, id).GetError()
	return &spec, err
}

func (r *ItemsSpecificationsRepository) FindAll() ([]models.ItemSpecification, error) {
	var specs []models.ItemSpecification
	err := r.db.Find(&specs).GetError()
	return specs, err
}

func (r *ItemsSpecificationsRepository) FindByItemID(itemID uint) ([]models.ItemSpecification, error) {
	var specs []models.ItemSpecification
	err := r.db.Where("item_id = ?", itemID).Find(&specs).GetError()
	return specs, err
}

func (r *ItemsSpecificationsRepository) Update(spec *models.ItemSpecification) error {
	return r.db.Save(spec).GetError()
}

func (r *ItemsSpecificationsRepository) Delete(id uint) error {
	return r.db.Delete(&models.ItemSpecification{}, id).GetError()
}
