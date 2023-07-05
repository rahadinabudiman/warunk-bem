package dtos

type InsertTransaksiResponse struct {
	Name       string `json:"name"`
	ProdukName string `json:"produk_name"`
	Total      int64  `json:"total"`
}
