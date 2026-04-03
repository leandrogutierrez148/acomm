package outbound

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
	Name     string    `json:"name"`
	Price    float64   `json:"price"`
	Category string    `json:"category"`
	Images   []string  `json:"images,omitempty"`
	Variants []Variant `json:"variants,omitempty"`
}

type Variant struct {
	Name   string   `json:"name"`
	SKU    string   `json:"sku"`
	Price  float64  `json:"price"`
	Images []string `json:"images,omitempty"`
}

type Pagination struct {
	Offset     int   `json:"offset"`
	Limit      int   `json:"limit"`
	TotalCount int64 `json:"total_count"`
}
