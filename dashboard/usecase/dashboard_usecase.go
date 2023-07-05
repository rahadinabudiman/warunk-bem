package usecase

import (
	"context"
	"errors"
	"time"
	"warunk-bem/domain"
)

type dashboardUsecase struct {
	UserRepo       domain.UserRepository
	UserAmountRepo domain.UserAmountRepository
	ProdukRepo     domain.ProdukRepository
	TransaksiRepo  domain.TransaksiRepository
	contextTimeout time.Duration
}

func NewDashboardUsecase(UserRepo domain.UserRepository, UserAmountRepo domain.UserAmountRepository, ProdukRepo domain.ProdukRepository, TransaksiRepo domain.TransaksiRepository, contextTimeout time.Duration) domain.DashboardUsecase {
	return &dashboardUsecase{
		UserRepo:       UserRepo,
		UserAmountRepo: UserAmountRepo,
		ProdukRepo:     ProdukRepo,
		TransaksiRepo:  TransaksiRepo,
		contextTimeout: contextTimeout,
	}
}

func (du *dashboardUsecase) GetDashboardData(c context.Context, userID string, rp int64, p int64, filter interface{}, setsort interface{}) (*domain.DashboardData, error) {
	ctx, cancel := context.WithTimeout(c, du.contextTimeout)
	defer cancel()

	// Mengambil saldo pengguna
	saldo, err := du.UserAmountRepo.FindOne(ctx, userID)
	if err != nil {
		return nil, errors.New("failed to get user's balance")
	}

	// Mengambil profil pengguna
	profil, err := du.UserRepo.FindOne(ctx, userID)
	if err != nil {
		return nil, errors.New("failed to get user's profile")
	}

	// Mengambil daftar produk
	produkList, _, err := du.ProdukRepo.GetAllWithPage(ctx, rp, p, filter, setsort)
	if err != nil {
		return nil, errors.New("failed to get product list")
	}

	// Membentuk response data dashboard
	dashboardData := &domain.DashboardData{
		Saldo:  saldo,
		Profil: profil,
		Produk: produkList,
	}

	return dashboardData, nil
}
