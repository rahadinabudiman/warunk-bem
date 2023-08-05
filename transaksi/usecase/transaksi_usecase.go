package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"warunk-bem/domain"
	"warunk-bem/dtos"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransaksiUsecase struct {
	TransaksiRepo  domain.TransaksiRepository
	KeranjangRepo  domain.KeranjangRepository
	ProdukRepo     domain.ProdukRepository
	UserRepo       domain.UserRepository
	UserAmountRepo domain.UserAmountRepository
	contextTimeout time.Duration
}

func NewTransaksiUsecase(TransaksiRepo domain.TransaksiRepository, KeranjangRepo domain.KeranjangRepository, ProdukRepo domain.ProdukRepository, UserRepo domain.UserRepository, UserAmountRepo domain.UserAmountRepository, contextTimeout time.Duration) domain.TransaksiUsecase {
	return &TransaksiUsecase{
		TransaksiRepo:  TransaksiRepo,
		KeranjangRepo:  KeranjangRepo,
		ProdukRepo:     ProdukRepo,
		UserRepo:       UserRepo,
		UserAmountRepo: UserAmountRepo,
		contextTimeout: contextTimeout,
	}
}

// GetTransaksi godoc
// @Summary      Get Transaksi by UserID
// @Description  Get Transaksi by UserID
// @Tags         User - Transaksi
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.TransaksiAllByUserIDResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /transaksi [get]
// @Security BearerAuth
func (tu *TransaksiUsecase) FindAll(c context.Context, id string) (res []*dtos.RiwayatTransaksiResponse, err error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	transaksis, err := tu.TransaksiRepo.FindAllByUserId(ctx, id)
	if err != nil {
		return nil, errors.New("cannot get user transaksi")
	}

	for _, transaksi := range transaksis {
		produk, err := tu.ProdukRepo.FindOne(ctx, transaksi.ProdukID.Hex())
		if err != nil {
			return nil, errors.New("cannot get produk")
		}

		totalharga := produk.Price * transaksi.Total

		riwayatTransaksi := &dtos.RiwayatTransaksiResponse{
			Name:      produk.Name,
			CreatedAt: transaksi.CreatedAt.Format("2006-01-02"),
			Waktu:     transaksi.CreatedAt.Format("15:04:05"),
			Harga:     totalharga,
			Total:     transaksi.Total,
			Image:     produk.Image,
		}

		res = append(res, riwayatTransaksi)
	}

	return res, nil
}

// AddTransaction godoc
// @Summary      Add Transaksi
// @Description  Add Transaksi
// @Tags         User - Transaksi
// @Accept       json
// @Produce      json
// @Param        request body dtos.TransaksiRequest true "Payload Body [RAW]"
// @Success      201 {object} dtos.TransaksiCreatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /transaksi [post]
// @Security BearerAuth
func (tu *TransaksiUsecase) InsertOne(ctx context.Context, req *dtos.InsertTransaksiRequest) (*dtos.InsertTransaksiResponse, error) {
	var res *dtos.InsertTransaksiResponse

	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	produk, err := tu.ProdukRepo.FindOne(ctx, req.ProdukID.Hex())
	if err != nil {
		return res, err
	}

	user, err := tu.UserRepo.FindOne(ctx, req.UserID.Hex())
	if err != nil {
		return res, err
	}

	if req.Total == 0 {
		return nil, errors.New("harap masukan jumlah produk yang ingin dibeli")
	}

	if produk.Stock == 0 {
		return nil, errors.New("produk telah habis terjual")
	}

	if produk.Stock < int64(req.Total) {
		return nil, errors.New("stok produk tidak mencukupi")
	}

	hargaProduk := produk.Price
	TotalBelanja := hargaProduk * int64(req.Total)

	saldo, err := tu.UserAmountRepo.FindOne(ctx, req.UserID.Hex())
	if err != nil {
		return res, errors.New("cannot get useramount")
	}

	if saldo.Amount < float64(TotalBelanja) {
		return nil, errors.New("saldo tidak mencukupi")
	}

	saldoAkhir := saldo.Amount - float64(TotalBelanja)
	saldo.Amount = saldoAkhir

	_, err = tu.UserAmountRepo.UpdateOne(ctx, saldo, saldo.ID.Hex())
	if err != nil {
		return res, err
	}

	produk.Stock = produk.Stock - int64(req.Total)
	produk.UpdatedAt = time.Now()

	_, err = tu.ProdukRepo.UpdateOne(ctx, produk, produk.ID.Hex())
	if err != nil {
		return res, err
	}

	req.ID = primitive.NewObjectID()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	transaksireq := &domain.Transaksi{
		ID:        req.ID,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
		UserID:    req.UserID,
		ProdukID:  req.ProdukID,
		Total:     int64(req.Total),
		Status:    "Berhasil",
	}

	resp, err := tu.TransaksiRepo.InsertOne(ctx, transaksireq)
	if err != nil {
		return res, err
	}

	res = &dtos.InsertTransaksiResponse{
		Name:       user.Name,
		ProdukName: produk.Name,
		Total:      resp.Total,
	}

	return res, nil
}

// AddTransaction godoc
// @Summary      Add Transaksi By Keranjang
// @Description  Add Transaksi By Keranjang
// @Tags         User - Transaksi
// @Accept       json
// @Produce      json
// @Param        request body dtos.InsertTransaksiKeranjangRequest true "Payload Body [RAW]"
// @Success      201 {object} dtos.TransaksiCreatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /transaksi/keranjang [post]
// @Security BearerAuth
func (tu *TransaksiUsecase) InsertByKeranjang(ctx context.Context, req *dtos.InsertTransaksiKeranjangRequest) (*dtos.InsertTransaksiResponse, error) {
	var res *dtos.InsertTransaksiResponse

	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	// Dapatkan data pengguna berdasarkan ID
	user, err := tu.UserRepo.FindOne(ctx, req.UserID.Hex())
	if err != nil {
		return res, errors.New("cannot get user")
	}

	// Dapatkan data keranjang berdasarkan ID
	keranjang, err := tu.KeranjangRepo.FindOneKeranjang(ctx, req.ID.Hex())
	if err != nil {
		return res, errors.New("cannot get keranjang")
	}

	// Periksa saldo pengguna:
	saldo, err := tu.UserAmountRepo.FindOne(ctx, req.UserID.Hex())
	if err != nil {
		return res, errors.New("cannot get user amount")
	}

	// Iterasi setiap produk dalam keranjang
	for _, produk := range keranjang.Produk {
		// Dapatkan data produk berdasarkan ID produk
		p, err := tu.ProdukRepo.FindOne(ctx, produk.ID.Hex())
		if err != nil {
			return res, err
		}

		// Periksa stok produk
		if p.Stock == 0 {
			return nil, fmt.Errorf("produk '%s' telah habis terjual", p.Name)
		}

		if p.Stock < produk.Stock {
			return nil, fmt.Errorf("stok produk '%s' tidak mencukupi", p.Name)
		}

		// Hitung total belanja untuk produk ini
		hargaProduk := p.Price
		totalBelanja := hargaProduk * int64(produk.Stock)

		// Periksa saldo pengguna
		if saldo.Amount < float64(totalBelanja) {
			return nil, fmt.Errorf("saldo tidak mencukupi untuk membeli produk '%s'", p.Name)
		}

		// Update saldo pengguna
		saldo.Amount -= float64(totalBelanja)
		_, err = tu.UserAmountRepo.UpdateOne(ctx, saldo, saldo.ID.Hex())
		if err != nil {
			return res, errors.New("cannot update saldo user")
		}

		// Update stok produk
		p.Stock -= produk.Stock
		p.UpdatedAt = time.Now()
		_, err = tu.ProdukRepo.UpdateOne(ctx, p, p.ID.Hex())
		if err != nil {
			return res, errors.New("cannot update produk stock")
		}

		// Insert transaksi
		transaksi := &domain.Transaksi{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			ProdukID:  p.ID,
			Total:     produk.Stock,
			Status:    "Berhasil",
		}

		_, err = tu.TransaksiRepo.InsertOne(ctx, transaksi)
		if err != nil {
			return res, errors.New("transaksi gagal")
		}
	}

	// Delete Keranjang Jika Transaksi Berhasil
	err = tu.KeranjangRepo.DeleteOne(ctx, keranjang.ID.Hex())
	if err != nil {
		return res, errors.New("cannot delete keranjang")
	}

	// Buat respons transaksi
	res = &dtos.InsertTransaksiResponse{
		Name:       user.Name,
		ProdukName: keranjang.Produk[0].Name, // Ambil nama produk pertama dalam keranjang
		Total:      int64(keranjang.Total),
	}

	return res, nil
}
