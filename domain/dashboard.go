package domain

import (
	"context"
)

type DashboardData struct {
	Saldo  *UserAmount `json:"saldo"`
	Profil *User       `json:"profil"`
	Produk []Produk    `json:"produk"`
}

type DashboardRepository interface {
	DashboardGetAll() ([]Produk, *User, error)
}

type DashboardUsecase interface {
	GetDashboardData(c context.Context, userID string, rp int64, p int64, filter interface{}, setsort interface{}) (*DashboardData, error)
}
