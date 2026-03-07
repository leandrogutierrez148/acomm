package repositories

import (
	"github.com/lgutierrez148/acomm/internal/interfaces"
	"github.com/lgutierrez148/acomm/internal/models"
)

type ItemsRepository struct {
	db interfaces.IDatabase
}

func NewItemsRepository(db interfaces.IDatabase) *ItemsRepository {
	return &ItemsRepository{db: db}
}

func (r *ItemsRepository) Create(item *models.Item) error {
	return r.db.Create(item).GetError()
}

func (r *ItemsRepository) FindByID(id uint) (*models.Item, error) {
	var item models.Item
	err := r.db.First(&item, id).GetError()
	return &item, err
}

func (r *ItemsRepository) FindAll() ([]models.Item, error) {
	var items []models.Item
	err := r.db.Find(&items).GetError()
	return items, err
}

func (r *ItemsRepository) FindByProductID(productID uint) ([]models.Item, error) {
	var items []models.Item
	err := r.db.Where("product_id = ?", productID).Find(&items).GetError()
	return items, err
}

func (r *ItemsRepository) Update(item *models.Item) error {
	return r.db.Save(item).GetError()
}

func (r *ItemsRepository) Delete(id uint) error {
	return r.db.Delete(&models.Item{}, id).GetError()
}
