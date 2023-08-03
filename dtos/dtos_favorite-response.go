package dtos


import "time"


type DeleteFavoriteRequest struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
}


type InsertFavoriteRequest struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    string    `json:"user_id"`
	ProdukID  string    `json:"produk_id"`
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
