package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterUserRequest struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	Name            string             `json:"name" form:"nama" validate:"required" example:"Rahadina Budiman Sundara"`
	Username        string             `json:"username" form:"username" validate:"required" example:"r4ha"`
	Email           string             `json:"email" form:"email" validate:"required,email" example:"me@r4ha.com"`
	Password        string             `json:"password" form:"password" validate:"gte=6" example:"rahadinabudimansundara"`
	PasswordConfirm string
	Verified        bool   `json:"verified" form:"verified" example:"False"`
	Role            string `json:"role" form:"role" gorm:"type:enum('Admin', 'Super Admin');default:'Admin'; not-null" example:"Admin"`
}

type RegisterUserAmountRequest struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Amount    float64            `bson:"amount" json:"amount"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" form:"nama" validate:"required" example:"Rahadina Budiman Sundara"`
	Username string `json:"username" form:"username" validate:"required" example:"r4ha"`
	Email    string `json:"email" form:"email" validate:"required,email" example:"me@r4ha.com"`
}

type TopUpSaldoRequest struct {
	Email  string  `json:"email" form:"email"`
	Amount float64 `json:"amount" form:"amount" validate:"required" example:"100000"`
}

type DeleteUserRequest struct {
	Password string `json:"password" form:"password" validate:"required" example:"rahadinabudimansundara"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
