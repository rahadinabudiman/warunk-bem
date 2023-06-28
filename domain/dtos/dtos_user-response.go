package dtos

import "go.mongodb.org/mongo-driver/bson/primitive"

type RegisterUserResponse struct {
	Name     string `json:"name" example:"Rahadina Budiman Sundara"`
	Username string `json:"username" example:"r4ha"`
	Email    string `json:"email" example:"r4ha@proton.me"`
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

type LoginUserResponse struct {
	Username string `json:"username" validate:"required"`
	Token    string `json:"token" form:"token"`
}
