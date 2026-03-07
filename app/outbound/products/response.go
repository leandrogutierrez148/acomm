package products

type GetProductsResponse struct {
	Products []Product `json:"products"`
}

type GetProductsPagedResponse struct {
	Products   []Product  `json:"products"`
	Pagination Pagination `json:"pagination"`
}

type GetProductResponse struct {
	Product Product `json:"product"`
}

type Product struct {
	Code     string    `json:"code"`
	Price    float64   `json:"price"`
	Category string    `json:"category"`
	Variants []Variant `json:"variants,omitempty"`
}

type Variant struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

type Pagination struct {
	Offset     int   `json:"offset"`
	Limit      int   `json:"limit"`
	TotalCount int64 `json:"total_count"`
}
