package orders

import "time"

type GetOrdersResponse struct {
	Orders []Order `json:"orders"`
}

type GetOrderResponse struct {
	Order Order `json:"order"`
}

type CreateOrderResponse struct {
	Order Order `json:"order"`
}

type UpdateOrderResponse struct {
	Order Order `json:"order"`
}

type Order struct {
	ID              int         `json:"id"`
	CustomerEmail   string      `json:"customer_email"`
	CustomerName    string      `json:"customer_name"`
	CustomerAddress string      `json:"customer_address"`
	CustomerPhone   string      `json:"customer_phone"`
	Status          string      `json:"status"`
	Items           []ItemOrder `json:"items,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type ItemOrder struct {
	ItemID  int `json:"item_id"`
	OrderID int `json:"order_id"`
}
