package usecase

import (
	"context"
	"time"
	"warunk-bem/domain"
	"warunk-bem/domain/dtos"
	"warunk-bem/user/usecase/helpers"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(u domain.UserRepository, to time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       u,
		contextTimeout: to,
	}
}

func (u *userUsecase) InsertOne(c context.Context, req *dtos.RegisterUserRequest) (*dtos.RegisterUserResponse, error) {
	var res *dtos.RegisterUserResponse

	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	req.ID = primitive.NewObjectID()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	passwordHash, err := helpers.HashPassword(req.Password)
	if err != nil {
		return res, err
	}

	req.Password = passwordHash
	req.Verified = false

	CreateUser := &domain.User{
		ID:        req.ID,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
		Name:      req.Name,
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		Verified:  req.Verified,
		Role:      req.Role,
	}

	createdUser, err := u.userRepo.InsertOne(ctx, CreateUser)
	if err != nil {
		return res, err
	}

	res = &dtos.RegisterUserResponse{
		ID:       createdUser.ID.Hex(),
		Name:     createdUser.Name,
		Email:    createdUser.Email,
		Username: createdUser.Username,
		Verified: createdUser.Verified,
		Role:     createdUser.Role,
	}

	return res, nil
}

func (u *userUsecase) FindOne(c context.Context, id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err := u.userRepo.FindOne(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (u *userUsecase) GetAllWithPage(c context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]domain.User, int64, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, count, err := u.userRepo.GetAllWithPage(ctx, rp, p, filter, setsort)
	if err != nil {
		return res, count, err
	}

	return res, count, nil
}

func (u *userUsecase) UpdateOne(c context.Context, m *domain.User, id string) (*domain.User, error) {

	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err := u.userRepo.UpdateOne(ctx, m, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (u *userUsecase) DeleteOne(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	err := u.userRepo.DeleteOne(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
