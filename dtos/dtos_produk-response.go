package dtos

type InsertProdukResponse struct {
	Name     string `bson:"name" json:"name"`
	Slug     string `bson:"slug" json:"slug"`
	Detail   string `bson:"detail" json:"detail"`
	Price    int64  `bson:"price" json:"price"`
	Stock    int64  `bson:"stock" json:"stock"`
	Category string `bson:"category" json:"category"`
	Image    string `bson:"image" json:"image"`
}

type ProdukDetailResponse struct {
	ID       string `bson:"_id" json:"id"`
	Name     string `bson:"name" json:"name"`
	Slug     string `bson:"slug" json:"slug"`
	Detail   string `bson:"detail" json:"detail"`
	Price    int64  `bson:"price" json:"price"`
	Stock    int64  `bson:"stock" json:"stock"`
	Category string `bson:"category" json:"category"`
	Image    string `bson:"image" json:"image"`
}

type GetAllProdukResponse struct {
	Total       int64                   `json:"total"`
	PerPage     int64                   `json:"per_page"`
	CurrentPage int64                   `json:"current_page"`
	LastPage    int64                   `json:"last_page"`
	From        int64                   `json:"from"`
	To          int64                   `json:"to"`
	Produk      []*ProdukDetailResponse `json:"produks"`
}
