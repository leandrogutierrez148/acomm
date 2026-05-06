package outbound

type GetCategoriesResponse struct {
	Categories []Category `json:"categories"`
}

type CreateCategoryResponse struct {
	Message string `json:"message"`
}

type Category struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
