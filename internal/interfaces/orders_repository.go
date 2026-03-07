package interfaces

import "github.com/lgutierrez148/acomm/internal/models"

type IOrdersRepository interface {
	Create(order *models.Order) error
	FindByID(id int) (*models.Order, error)
	FindAll() ([]models.Order, error)
	FindByCustomerEmail(email string) ([]models.Order, error)
	Update(order *models.Order) error
	Delete(id int) error
}
