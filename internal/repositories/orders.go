package repositories

import (
	"github.com/lgutierrez148/acomm/internal/interfaces"
	"github.com/lgutierrez148/acomm/internal/models"
)

type OrdersRepository struct {
	db interfaces.IDatabase
}

func NewOrdersRepository(db interfaces.IDatabase) *OrdersRepository {
	return &OrdersRepository{db: db}
}

func (r *OrdersRepository) Create(order *models.Order) error {
	return r.db.Create(order).GetError()
}

func (r *OrdersRepository) FindByID(id int) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items").First(&order, id).GetError()
	return &order, err
}

func (r *OrdersRepository) FindAll() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Items").Find(&orders).GetError()
	return orders, err
}

func (r *OrdersRepository) FindByCustomerEmail(email string) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Items").Where("customer_email = ?", email).Find(&orders).GetError()
	return orders, err
}

func (r *OrdersRepository) Update(order *models.Order) error {
	return r.db.Save(order).GetError()
}

func (r *OrdersRepository) Delete(id int) error {
	return r.db.Delete(&models.Order{}, id).GetError()
}
