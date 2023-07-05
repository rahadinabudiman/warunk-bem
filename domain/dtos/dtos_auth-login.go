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
