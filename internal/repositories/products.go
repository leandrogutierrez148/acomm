package repositories

import (
	"strconv"

	"github.com/lgutierrez148/acomm/internal/interfaces"
	"github.com/lgutierrez148/acomm/internal/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductsRepository struct {
	db interfaces.IDatabase
}

func NewProductsRepository(db interfaces.IDatabase) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := r.db.Preload("Category").Preload("Items").Preload("Brand").Find(&products).GetError(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductsRepository) GetProductsPaginated(offset, limit int) ([]models.Product, int64, error) {
	var products []models.Product

	if err := r.db.Preload("Category").Preload("Items").Preload("Brand").Offset(offset).Limit(limit).Find(&products).GetError(); err != nil {
		return nil, 0, err
	}

	var total int64
	r.db.Model(&models.Product{}).Count(&total)

	return products, total, nil
}

func (r *ProductsRepository) SearchProductsPaginated(offset, limit int, category string, maxPrice decimal.Decimal) ([]models.Product, int64, error) {
	var products []models.Product

	query := r.db.Model(&models.Product{}).
		Preload("Category").
		Preload("Items").
		Preload("Brand").
		Joins("JOIN product_categories ON product_categories.id = products.category_id")

	if maxPrice.Compare(decimal.Zero) > 0 {
		query = query.Joins("LEFT JOIN items ON items.product_id = products.id").
			Where("items.price <= ?", maxPrice)
	}

	if category != "" {
		query = query.Where("product_categories.name = ?", category)
	}

	query = query.Group("products.id")

	if err := query.Offset(offset).Limit(limit).Find(&products).GetError(); err != nil {
		return nil, 0, err
	}

	var total int64
	countQuery := r.db.Model(&models.Product{}).
		Joins("JOIN product_categories ON product_categories.id = products.category_id")

	if maxPrice.Compare(decimal.Zero) > 0 {
		countQuery = countQuery.Joins("LEFT JOIN items ON items.product_id = products.id").
			Where("items.price <= ?", maxPrice)
	}

	if category != "" {
		countQuery = countQuery.Where("product_categories.name = ?", category)
	}

	countQuery.Distinct("products.id").Count(&total)

	return products, total, nil
}

func (r *ProductsRepository) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Category").Preload("Items").Preload("Brand").Where("products.id = ?", id).First(&product).GetError(); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductsRepository) GetProductByCode(code string) (*models.Product, error) {
	id, err := strconv.ParseUint(code, 10, 32)
	if err != nil {
		return nil, nil
	}
	return r.GetProductByID(uint(id))
}
