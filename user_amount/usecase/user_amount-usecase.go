package usecase

import (
	"context"
	"errors"
	"time"
	"warunk-bem/domain"
	"warunk-bem/dtos"
)

type UserAmountUsecase struct {
	UserAmountRepo domain.UserAmountRepository
	UserRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewUserAmountUsecase(ua domain.UserAmountRepository, u domain.UserRepository, timeout time.Duration) domain.UserAmountUsecase {
	return &UserAmountUsecase{
		UserAmountRepo: ua,
		UserRepo:       u,
		contextTimeout: timeout,
	}
}

func (uas *UserAmountUsecase) TopUpSaldo(ctx context.Context, req *dtos.TopUpSaldoRequest) (res *dtos.TopUpSaldoResponse, err error) {
	var (
		userAmount *domain.UserAmount
		user       *domain.User
	)

	ctx, cancel := context.WithTimeout(ctx, uas.contextTimeout)
	defer cancel()

	user, err = uas.UserRepo.FindEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("email tidak ditemukan")
	}

	userAmount, err = uas.UserAmountRepo.FindOne(ctx, user.ID.Hex())
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	userAmount.Amount += req.Amount

	_, err = uas.UserAmountRepo.UpdateOne(ctx, userAmount, userAmount.ID.Hex())
	if err != nil {
		return nil, errors.New("tidak dapat menambahkan saldo")
	}

	Message := "Saldo berhasil ditambahkan ke akun"

	res = &dtos.TopUpSaldoResponse{
		Name:    user.Name,
		Amount:  req.Amount,
		Message: Message,
	}

	return res, nil
}
