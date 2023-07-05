package http

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"warunk-bem/domain"
	"warunk-bem/domain/dtos"
	"warunk-bem/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	UsrUsecase domain.UserUsecase
}

func NewUserHandler(router *gin.RouterGroup, protected *gin.RouterGroup, protectedAdmin *gin.RouterGroup, uu domain.UserUsecase) {
	handler := &UserHandler{
		UsrUsecase: uu,
	}

	// Main API
	api := router.Group("/user")
	protected = protected.Group("/user")
	protectedAdmin = protectedAdmin.Group("/user")

	api.POST("", handler.InsertOne)
	api.POST("/activation", handler.VerifyAccount)
	protected.POST("/verify", handler.VerifyLogin)
	protected.GET("/:id", handler.FindOne)
	protectedAdmin.GET("", handler.GetAll)
	protected.PUT("/:id", handler.UpdateOne)
	protected.DELETE("/:id", handler.DeleteOne)
}

func isRequestValid(m *dtos.RegisterUserRequest) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (user *UserHandler) InsertOne(c *gin.Context) {
	var (
		usr dtos.RegisterUserRequest
		err error
	)

	err = c.ShouldBindJSON(&usr)
	if err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dtos.NewErrorResponse(
				http.StatusUnprocessableEntity,
				"Filed Cannot Be Empty",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	var ok bool
	if ok, err = isRequestValid(&usr); !ok {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Bad Request",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := user.UsrUsecase.InsertOne(ctx, &usr)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			dtos.NewErrorResponse(
				http.StatusInternalServerError,
				"Cannot Insert Data",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	message := "We sent an email with a verification code to " + result.Email
	c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Success",
			message,
		),
	)
}

func (user *UserHandler) FindOne(c *gin.Context) {
	dataUser, err := middlewares.IsUser(c)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Please login first",
				dtos.GetErrorData(err),
			),
		)
		return
	}
	id := c.Param("id")

	if dataUser != id {
		c.JSON(
			http.StatusForbidden,
			dtos.NewResponseMessage(
				http.StatusForbidden,
				"Forbidden",
			),
		)
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := user.UsrUsecase.FindOne(ctx, id)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			dtos.NewErrorResponse(
				http.StatusInternalServerError,
				"Cannot Find Data",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Success",
			result,
		),
	)
}

func (user *UserHandler) GetAll(c *gin.Context) {
	var (
		res   []dtos.UserProfileResponse
		count int64
	)

	rp, err := strconv.ParseInt(c.Query("rp"), 10, 64)
	if err != nil {
		rp = 25
	}

	page, err := strconv.ParseInt(c.Query("p"), 10, 64)
	if err != nil {
		page = 1
	}

	filters := bson.D{{Key: "name", Value: primitive.Regex{Pattern: ".*" + c.Query("name") + ".*", Options: "i"}}}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, count, err = user.UsrUsecase.GetAllWithPage(ctx, rp, page, filters, nil)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			dtos.NewErrorResponse(
				http.StatusInternalServerError,
				"Cannot Find Data",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	result := dtos.GetAllUserResponse{
		Total:       count,
		PerPage:     rp,
		CurrentPage: page,
		LastPage:    int64(math.Ceil(float64(count) / float64(rp))),
		From:        page*rp - rp + 1,
		To:          page * rp,
		User:        res,
	}

	c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Success",
			result,
		),
	)
}

func (user *UserHandler) UpdateOne(c *gin.Context) {
	id := c.Param("id")

	var (
		usr dtos.UpdateUserRequest
		err error
	)

	err = c.ShouldBindJSON(&usr)
	if err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dtos.NewErrorResponse(
				http.StatusUnprocessableEntity,
				"Filed Cannot Be Empty",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := user.UsrUsecase.UpdateOne(ctx, &usr, id)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			dtos.NewErrorResponse(
				http.StatusInternalServerError,
				"Cannot Update Data",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Success",
			result,
		),
	)
}

func (user *UserHandler) DeleteOne(c *gin.Context) {
	id := c.Param("id")

	idUser, err := middlewares.IsUser(c)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Please login first before access to this",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	if idUser != id {
		c.JSON(
			http.StatusForbidden,
			dtos.NewResponseMessage(
				http.StatusForbidden,
				"Forbidden",
			),
		)
		return
	}

	var req dtos.DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dtos.NewErrorResponse(
				http.StatusUnprocessableEntity,
				"Please verify your password and try again",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	_, err = user.UsrUsecase.DeleteOne(ctx, id, req)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			dtos.NewErrorResponse(
				http.StatusInternalServerError,
				"Cannot Delete Data",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dtos.NewResponseMessage(
			http.StatusOK,
			"Success",
		),
	)
}

func (user *UserHandler) VerifyLogin(c *gin.Context) {
	var (
		req dtos.VerifyLoginRequest
		err error
	)

	_, err = middlewares.IsUser(c)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Please login first before access to this",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dtos.NewErrorResponse(
				http.StatusUnprocessableEntity,
				"Filed Cannot Be Empty",
				dtos.GetErrorData(err),
			),
		)
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := user.UsrUsecase.VerifyLogin(ctx, req.Code)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Verification Code is Wrong",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Success",
			result,
		),
	)
}

func (user *UserHandler) VerifyAccount(c *gin.Context) {
	var (
		req dtos.ActivationAccountRequest
		err error
	)

	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dtos.NewErrorResponse(
				http.StatusUnprocessableEntity,
				"Filed Cannot Be Empty",
				dtos.GetErrorData(err),
			),
		)
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := user.UsrUsecase.VerifyAccount(ctx, req.Code)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Verification Code is Wrong",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Your Activation has success",
			res,
		),
	)
}
