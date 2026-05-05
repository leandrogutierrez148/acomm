package outbound

type GetBrandsResponse struct {
	Brands []Brand `json:"brands"`
}

type GetBrandResponse struct {
	Brand Brand `json:"brand"`
}

type CreateBrandResponse struct {
	Brand Brand `json:"brand"`
}

type UpdateBrandResponse struct {
	Brand Brand `json:"brand"`
}

type Brand struct {
	ID                 uint   `json:"id"`
	Name               string `json:"name"`
	Title              string `json:"title"`
	ImageURL           string `json:"image_url"`
	MetaTagDescription string `json:"meta_tag_description"`
	IsActive           bool   `json:"is_active"`
}
