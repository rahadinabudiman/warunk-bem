package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"warunk-bem/domain"

	"github.com/go-redis/redis/v8"
)

type dashboardUsecase struct {
	DashboardRepo  domain.DashboardRepository
	UserRepo       domain.UserRepository
	UserAmountRepo domain.UserAmountRepository
	ProdukRepo     domain.ProdukRepository
	TransaksiRepo  domain.TransaksiRepository
	RedisClient    *redis.Client
	contextTimeout time.Duration
}

func NewDashboardUsecase(DashboardRepo domain.DashboardRepository, UserRepo domain.UserRepository, UserAmountRepo domain.UserAmountRepository, ProdukRepo domain.ProdukRepository, TransaksiRepo domain.TransaksiRepository, RedisClient *redis.Client, contextTimeout time.Duration) domain.DashboardUsecase {
	return &dashboardUsecase{
		DashboardRepo:  DashboardRepo,
		UserRepo:       UserRepo,
		UserAmountRepo: UserAmountRepo,
		ProdukRepo:     ProdukRepo,
		TransaksiRepo:  TransaksiRepo,
		RedisClient:    RedisClient,
		contextTimeout: contextTimeout,
	}
}

func (du *dashboardUsecase) GetDashboardData(c context.Context, userID string, rp int64, p int64, filter interface{}, setsort interface{}) (*domain.DashboardData, error) {

	var res *domain.DashboardData

	cachekey := "dashboard_data:" + userID
	val, err := du.RedisClient.Get(c, cachekey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(val), &res); err != nil {
			return nil, errors.New("failed to unmarshal cache")
		}
		return res, nil
	}

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

	cacheValue, err := json.Marshal(dashboardData)
	if err == nil {
		du.RedisClient.Set(c, cachekey, cacheValue, 10*time.Minute)
	}

	return dashboardData, nil
}
