package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterUserResponse struct {
	Name     string `json:"name" example:"Rahadina Budiman Sundara"`
	Username string `json:"username" example:"r4ha"`
	Email    string `json:"email" example:"r4ha@proton.me"`
}

type RegisterUserResponseVerification struct {
	Email   string `json:"email" example:"r4ha@proton.me"`
	Message string `json:"message" form:"message"`
}

type RegisterUserAmountResponse struct {
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	Amount float64            `bson:"amount" json:"amount"`
}

type UserProfileResponse struct {
	Name     string `json:"name" example:"Rahadina Budiman Sundara"`
	Username string `json:"username" example:"r4ha"`
	Email    string `json:"email" example:"r4ha@proton.me"`
}

type UserDetailResponse struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	Email     string             `bson:"email" json:"email" validate:"required"`
	Username  string             `bson:"username" json:"username" validate:"required"`
	Role      string             `bson:"role" json:"role" validate:"required"`
}

type GetAllUserResponse struct {
	Total       int64                 `json:"total"`
	PerPage     int64                 `json:"per_page"`
	CurrentPage int64                 `json:"current_page"`
	LastPage    int64                 `json:"last_page"`
	From        int64                 `json:"from"`
	To          int64                 `json:"to"`
	User        []UserProfileResponse `json:"users"`
}

type UpdateUserResponse struct {
	Name     string `json:"name" form:"nama" validate:"required" example:"Rahadina Budiman Sundara"`
	Username string `json:"username" form:"username" validate:"required" example:"r4ha"`
	Email    string `json:"email" form:"email" validate:"required,email" example:"me@r4ha.com"`
}

type TopUpSaldoResponse struct {
	Name    string  `json:"name" form:"name"`
	Amount  float64 `json:"amount" form:"amount" validate:"required" example:"100000"`
	Message string  `json:"message" form:"message"`
}

type LoginUserResponse struct {
	Username string `json:"username" validate:"required"`
	Token    string `json:"token" form:"token"`
	Message  string `json:"message" form:"message"`
}

type LogoutUserResponse struct {
	Message string `json:"message"`
}
