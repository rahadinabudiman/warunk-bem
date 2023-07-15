package dtos

type InsertTransaksiResponse struct {
	Name       string `json:"name"`
	ProdukName string `json:"produk_name"`
	Total      int64  `json:"total"`
}

type RiwayatTransaksiResponse struct {
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Waktu     string `json:"waktu"`
	Harga     int64  `json:"harga"`
	Total     int64  `json:"total"`
	Image     string `json:"image"`
}
