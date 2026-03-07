package repositories

import (
	"github.com/lgutierrez148/acomm/models"
	"gorm.io/gorm"
)

type OrdersRepository struct {
	db *gorm.DB
}

func NewOrdersRepository(db *gorm.DB) *OrdersRepository {
	return &OrdersRepository{db: db}
}

func (r *OrdersRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrdersRepository) FindByID(id int) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items").First(&order, id).Error
	return &order, err
}

func (r *OrdersRepository) FindAll() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Items").Find(&orders).Error
	return orders, err
}

func (r *OrdersRepository) FindByCustomerEmail(email string) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Items").Where("customer_email = ?", email).Find(&orders).Error
	return orders, err
}

func (r *OrdersRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *OrdersRepository) Delete(id int) error {
	return r.db.Delete(&models.Order{}, id).Error
}
