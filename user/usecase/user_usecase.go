package usecase

import (
	"context"
	"time"
	"warunk-bem/domain"

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

func (u *userUsecase) InsertOne(c context.Context, m *domain.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	res, err := u.userRepo.InsertOne(ctx, m)
	if err != nil {
		return res, err
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
