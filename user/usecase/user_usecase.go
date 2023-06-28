package usecase

import (
	"context"
	"errors"
	"strings"
	"time"
	"warunk-bem/domain"
	"warunk-bem/domain/dtos"
	"warunk-bem/user/usecase/helpers"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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

	err := helpers.ValidateUsername(req.Username)
	if err != nil {
		return res, errors.New(err.Error())
	}

	username, _ := u.userRepo.FindUsername(ctx, req.Username)
	if username.ID != [12]byte{} {
		return res, errors.New("username already exists")
	}

	req.Email = strings.ToLower(req.Email)
	email, _ := u.userRepo.FindEmail(ctx, req.Email)
	if email.ID != [12]byte{} {
		return res, errors.New("email already exist")
	}

	if req.Password != req.PasswordConfirm {
		return res, errors.New("password does not match")
	}

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
		Name:     createdUser.Name,
		Email:    createdUser.Email,
		Username: createdUser.Username,
	}

	return res, nil
}

func (u *userUsecase) FindOne(c context.Context, id string) (res *dtos.UserProfileResponse, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	req, err := u.userRepo.FindOne(ctx, id)
	if err != nil {
		return res, err
	}

	res = &dtos.UserProfileResponse{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
	}

	return res, nil
}

func (u *userUsecase) GetAllWithPage(c context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]dtos.UserProfileResponse, int64, error) {
	var res []dtos.UserProfileResponse

	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	req, count, err := u.userRepo.GetAllWithPage(ctx, rp, p, filter, setsort)
	if err != nil {
		return res, count, err
	}

	for _, v := range req {
		res = append(res, dtos.UserProfileResponse{
			Name:     v.Name,
			Username: v.Username,
			Email:    v.Email,
		})
	}

	return res, count, nil
}

func (u *userUsecase) UpdateOne(c context.Context, req *dtos.UpdateUserRequest, id string) (*dtos.UpdateUserResponse, error) {
	var (
		res *dtos.UpdateUserResponse
	)
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.userRepo.FindOne(ctx, id)
	if err != nil {
		return res, err
	}

	result.Name = req.Name
	result.Username = req.Username
	result.Email = req.Email

	resp, err := u.userRepo.UpdateOne(ctx, result, id)
	if err != nil {
		return res, err
	}

	res = &dtos.UpdateUserResponse{
		Name:     resp.Name,
		Username: resp.Username,
		Email:    resp.Email,
	}

	return res, nil
}

func (u *userUsecase) DeleteOne(c context.Context, id string, req dtos.DeleteUserRequest) (res dtos.ResponseMessage, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.FindOne(ctx, id)
	if err != nil {
		return res, errors.New("user not found")
	}

	err = helpers.ComparePassword(req.Password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return res, errors.New("password is incorrect")
	}

	err = u.userRepo.DeleteOne(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}
