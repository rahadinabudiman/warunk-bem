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

type RegisterUserResponse struct {
	ID       string `json:"id" example:"5f7b1a7d9b3b1a1b1a1b1a1b"`
	Name     string `json:"name" example:"Rahadina Budiman Sundara"`
	Username string `json:"username" example:"r4ha"`
	Email    string `json:"email" example:"r4ha@proton.me"`
	Verified bool   `json:"verified" example:"False"`
	Role     string `json:"role" example:"Admin"`
}
