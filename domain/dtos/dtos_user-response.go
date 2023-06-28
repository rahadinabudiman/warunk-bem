package dtos

type RegisterUserResponse struct {
	Name     string `json:"name" example:"Rahadina Budiman Sundara"`
	Username string `json:"username" example:"r4ha"`
	Email    string `json:"email" example:"r4ha@proton.me"`
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
