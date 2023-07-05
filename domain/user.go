package domain

import (
	"context"
	"time"
	"warunk-bem/dtos"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
	Name             string             `bson:"name" json:"name" validate:"required"`
	Email            string             `bson:"email" json:"email" validate:"required"`
	Username         string             `bson:"username" json:"username" validate:"required"`
	Password         string             `bson:"password" json:"password" validate:"required"`
	Verified         bool               `bson:"verified" json:"verified"`
	LoginVerif       int                `bson:"loginverif" json:"loginverif"`
	VerificationCode int                `bson:"verification" json:"verification"`
	ActivationCode   int                `bson:"activation_code" json:"activation_code"`
	Role             string             `bson:"role" json:"role" validate:"required"`
}

type UserRepository interface {
	InsertOne(ctx context.Context, req *User) (*User, error)
	FindOne(ctx context.Context, id string) (*User, error)
	FindUsername(ctx context.Context, username string) (*User, error)
	FindEmail(ctx context.Context, email string) (*User, error)
	FindVerificationCode(ctx context.Context, verification int) (*User, error)
	FindActivationCode(ctx context.Context, activation int) (*User, error)
	GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]User, int64, error)
	UpdateOne(ctx context.Context, user *User, id string) (*User, error)
	GetByCredential(ctx context.Context, req *dtos.LoginUserRequest) (*User, error)
	DeleteOne(ctx context.Context, id string) error
}
type UserUsecase interface {
	InsertOne(ctx context.Context, req *dtos.RegisterUserRequest) (*dtos.RegisterUserResponseVerification, error)
	FindOne(ctx context.Context, id string) (res *dtos.UserProfileResponse, err error)
	VerifyLogin(ctx context.Context, verification int) (res dtos.VerifyLoginResponse, err error)
	VerifyAccount(ctx context.Context, activation int) (res dtos.VerifyEmailResponse, err error)
	GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]dtos.UserProfileResponse, int64, error)
	UpdateOne(ctx context.Context, user *dtos.UpdateUserRequest, id string) (*dtos.UpdateUserResponse, error)
	DeleteOne(c context.Context, id string, req dtos.DeleteUserRequest) (res dtos.ResponseMessage, err error)
}
