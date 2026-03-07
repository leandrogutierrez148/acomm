package inbound

import "github.com/lgutierrez148/acomm/internal/models"

type ItemOrderRequest struct {
	ItemID int `json:"item_id"`
}

type CreateOrderRequest struct {
	CustomerEmail   string             `json:"customer_email"`
	CustomerName    string             `json:"customer_name"`
	CustomerAddress string             `json:"customer_address"`
	CustomerPhone   string             `json:"customer_phone"`
	Items           []ItemOrderRequest `json:"items"`
}

func (r *CreateOrderRequest) ToDomain() *models.Order {
	items := make([]models.ItemOrder, len(r.Items))
	for i, it := range r.Items {
		items[i] = models.ItemOrder{ItemID: it.ItemID}
	}
	return &models.Order{
		CustomerEmail:   r.CustomerEmail,
		CustomerName:    r.CustomerName,
		CustomerAddress: r.CustomerAddress,
		CustomerPhone:   r.CustomerPhone,
		Items:           items,
	}
}

type UpdateOrderRequest struct {
	CustomerEmail   string             `json:"customer_email"`
	CustomerName    string             `json:"customer_name"`
	CustomerAddress string             `json:"customer_address"`
	CustomerPhone   string             `json:"customer_phone"`
	Status          string             `json:"status"`
	Items           []ItemOrderRequest `json:"items,omitempty"`
}

func (r *UpdateOrderRequest) ToDomain() *models.Order {
	items := make([]models.ItemOrder, len(r.Items))
	for i, it := range r.Items {
		items[i] = models.ItemOrder{ItemID: it.ItemID}
	}
	return &models.Order{
		CustomerEmail:   r.CustomerEmail,
		CustomerName:    r.CustomerName,
		CustomerAddress: r.CustomerAddress,
		CustomerPhone:   r.CustomerPhone,
		Status:          r.Status,
		Items:           items,
	}
}
