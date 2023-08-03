package dtos

type StatusOKResponse struct {
	StatusCode int         `json:"status_code" example:"200"`
	Message    string      `json:"message" example:"Successfully"`
	Data       interface{} `json:"data"`
}
type LogoutAdminOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Message    string `json:"message" form:"message" example:"Logout Success"`
}

type VerifyEmailOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Username   string `json:"username" form:"username" validate:"required" example:"r4ha"`
	Message    string `json:"message" form:"message" example:"Email has been verified"`
}

type VerifyLoginOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Token      string `json:"token" form:"token" validate:"required" example:"a82jask2jafk1l2kam"`
	Message    string `json:"message" form:"message" example:"Email has been verified"`
}

type ForgotPasswordOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Email      string `json:"email" form:"email" example:"me@r4ha.com"`
	Message    string `json:"message" form:"message" example:"OTP has been sent to your email"`
}
type ChangePasswordOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Email      string `json:"email" form:"email" example:"me@r4ha.com"`
	Message    string `json:"message" form:"message" example:"Password has been reset successfully"`
}

type ChangePasswordAdminOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Message    string `json:"message" form:"message" example:"Password has been reset successfully"`
}

type ChangePasswordUserOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Message    string `json:"message" form:"message" example:"Password has been reset successfully"`
}

type LoginStatusOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Username   string `json:"username" form:"username" validate:"required" example:"r4ha"`
	Message    string `json:"message" form:"message" example:"Login Success"`
	Token      string `json:"token" form:"token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
}

type TopUpSaldoOKResponse struct {
	StatusCode int     `json:"status_code" example:"200"`
	Name       string  `json:"name" form:"name"`
	Amount     float64 `json:"amount" form:"amount" validate:"required" example:"100000"`
	Message    string  `json:"message" form:"message"`
}

type UserStatusOKResponse struct {
	StatusCode int                `json:"status_code" example:"200"`
	Message    string             `json:"message" example:"Successfully get user credentials"`
	Data       UserDetailResponse `json:"data"`
}

type UserCreeatedResponse struct {
	StatusCode int                `json:"status_code" example:"201"`
	Message    string             `json:"message" example:"Successfully registered"`
	Data       UserDetailResponse `json:"data"`
}

type ProdukCreatedResponse struct {
	StatusCode int                  `json:"status_code" example:"201"`
	Message    string               `json:"message" example:"Successfully registered"`
	Data       ProdukDetailResponse `json:"data"`
}

type FavoriteCreatedResponse struct {
	StatusCode int                    `json:"status_code" example:"201"`
	Message    string                 `json:"message" example:"Successfully registered"`
	Data       DetailFavoriteResponse `json:"data"`
}

type DeleteProductFavoriteResponse struct {
	StatusCode int                     `json:"status_code" example:"201"`
	Message    string                  `json:"message" example:"Successfully deleted"`
	Data       DelelteFavoriteResponse `json:"data"`
}

type ProdukOKResponse struct {
	StatusCode int                  `json:"status_code" example:"200"`
	Message    string               `json:"message" example:"Successfully"`
	Data       ProdukDetailResponse `json:"data"`
}

type GetAllUserResponses struct {
	StatusCode int                `json:"status_code" example:"201"`
	Message    string             `json:"message" example:"Successfully registered"`
	Data       UserDetailResponse `json:"data"`
}

type StatusOKDeletedResponse struct {
	StatusCode int         `json:"status_code" example:"200"`
	Message    string      `json:"message" example:"Successfully deleted"`
	Errors     interface{} `json:"errors"`
}

type BadRequestResponse struct {
	StatusCode int         `json:"status_code" example:"400"`
	Message    string      `json:"message" example:"Bad Request"`
	Errors     interface{} `json:"errors"`
}

type UnauthorizedResponse struct {
	StatusCode int         `json:"status_code" example:"401"`
	Message    string      `json:"message" example:"Unauthorized"`
	Errors     interface{} `json:"errors"`
}

type ForbiddenResponse struct {
	StatusCode int         `json:"status_code" example:"403"`
	Message    string      `json:"message" example:"Forbidden"`
	Errors     interface{} `json:"errors"`
}

type NotFoundResponse struct {
	StatusCode int         `json:"status_code" example:"404"`
	Message    string      `json:"message" example:"Not Found"`
	Errors     interface{} `json:"errors"`
}

type InternalServerErrorResponse struct {
	StatusCode int         `json:"status_code" example:"500"`
	Message    string      `json:"message" example:"Internal Server Error"`
	Errors     interface{} `json:"errors"`
}
