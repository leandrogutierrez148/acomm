package interfaces

import (
	"github.com/lgutierrez148/acomm/models"
	"github.com/shopspring/decimal"
)

type IProductsRepository interface {
	GetAllProducts() ([]models.Product, error)
	GetProductsPaginated(offset, limit int) ([]models.Product, int64, error)
	SearchProductsPaginated(offset, limit int, category string, maxPrice decimal.Decimal) ([]models.Product, int64, error)
	GetProductByID(id uint) (*models.Product, error)
	GetProductByCode(code string) (*models.Product, error)
}
