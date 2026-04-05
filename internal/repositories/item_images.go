package repositories

import (
	"github.com/leandrogutierrez148/acomm/internal/interfaces"
	"github.com/leandrogutierrez148/acomm/internal/models"
)

type ItemImagesRepository struct {
	db interfaces.IDatabase
}

func NewItemImagesRepository(db interfaces.IDatabase) *ItemImagesRepository {
	return &ItemImagesRepository{db: db}
}

func (r *ItemImagesRepository) Create(image *models.ItemImage) error {
	return r.db.Create(image).GetError()
}

func (r *ItemImagesRepository) FindByID(id uint) (*models.ItemImage, error) {
	var image models.ItemImage
	err := r.db.First(&image, id).GetError()
	return &image, err
}

func (r *ItemImagesRepository) FindAll() ([]models.ItemImage, error) {
	var images []models.ItemImage
	err := r.db.Find(&images).GetError()
	return images, err
}

func (r *ItemImagesRepository) FindByItemID(itemID uint) ([]models.ItemImage, error) {
	var images []models.ItemImage
	err := r.db.Where("item_id = ?", itemID).Find(&images).GetError()
	return images, err
}

func (r *ItemImagesRepository) Update(image *models.ItemImage) error {
	return r.db.Save(image).GetError()
}

func (r *ItemImagesRepository) Delete(id uint) error {
	return r.db.Delete(&models.ItemImage{}, id).GetError()
}
