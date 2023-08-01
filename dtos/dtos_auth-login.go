package dtos

type VerifyLoginRequest struct {
	Code int `json:"code" form:"code"`
}

type ActivationAccountRequest struct {
	Code int `json:"code" form:"code"`
}

type VerifyEmailResponse struct {
	Email   string `json:"email" form:"email" validate:"required" example:"r4ha"`
	Message string `json:"message" form:"message" example:"Email has been verified"`
}

type VerifyLoginResponse struct {
	Message string `json:"message" form:"message" example:"Email has been verified"`
	Token   string `json:"token" form:"token" example:"29eiekk10k3k3k"`
}
