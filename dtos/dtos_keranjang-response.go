package dtos

type InsertKeranjangResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Produk []struct {
		ProdukID string `json:"produk_id"`
		Name     string `bson:"name" json:"name"`
		Detail   string `bson:"detail" json:"detail"`
		Price    int64  `bson:"price" json:"price"`
		Stock    int64  `bson:"stock" json:"stock"`
		Category string `bson:"category" json:"category"`
		Image    string `bson:"image" json:"image" form:"image"`
		Total    int    `json:"total"`
	} `json:"produk"`
}
