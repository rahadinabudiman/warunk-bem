package dtos

type DetailKeranjang struct {
	ID     string               `bson:"_id" json:"id"`
	UserID string               `bson:"user_id" json:"user_id"`
	Produk ProdukDetailResponse `bson:"produk" json:"produk"`
	Total  int                  `bson:"total" json:"total"`
}

type InsertKeranjangRequest struct {
	ProdukID string `json:"produk_id"`
	Total    int    `json:"total"`
}

type UpdateKeranjangRequest struct {
	ID     string               `bson:"_id" json:"id"`
	Produk ProdukDetailResponse `json:"produk"`
}

type DeleteProductKeranjangRequest struct {
	KeranjangID string `json:"keranjang_id"`
	ProdukID    string `json:"produk_id"`
}

type DeleteProductKeranjang struct {
	Name string `json:"name"`
}
