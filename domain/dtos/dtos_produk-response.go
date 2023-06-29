package dtos

type InsertProdukResponse struct {
	Name     string `bson:"name" json:"name"`
	Detail   string `bson:"detail" json:"detail"`
	Price    int64  `bson:"price" json:"price"`
	Stock    int64  `bson:"stock" json:"stock"`
	Category string `bson:"category" json:"category"`
}

type ProdukDetailResponse struct {
	Name     string `bson:"name" json:"name"`
	Detail   string `bson:"detail" json:"detail"`
	Price    int64  `bson:"price" json:"price"`
	Stock    int64  `bson:"stock" json:"stock"`
	Category string `bson:"category" json:"category"`
}
