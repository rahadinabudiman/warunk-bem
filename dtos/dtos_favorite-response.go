package dtos

type DeleteFavoriteRequest struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
}

type DetailFavoriteResponse struct {
	ID     string                 `json:"id"`
	UserID string                 `json:"user_id"`
	Produk []ProdukDetailResponse `json:"produk"`
	Total  int                    `json:"total"`
}

type DelelteFavoriteResponse struct {
	Name string `json:"name"`
}
