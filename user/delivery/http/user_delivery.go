package http

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"warunk-bem/domain"
	"warunk-bem/domain/dtos"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	UsrUsecase domain.UserUsecase
}

func NewUserHandler(e *echo.Echo, uu domain.UserUsecase) {
	handler := &UserHandler{
		UsrUsecase: uu,
	}

	// Main API
	api := e.Group("/api/v1")
	user := api.Group("/user")

	user.POST("", handler.InsertOne)
	user.GET("/:id", handler.FindOne)
	user.GET("", handler.GetAll)
	user.PUT("/:id", handler.UpdateOne)
	user.DELETE("/:id", handler.DeleteOne)
}

func isRequestValid(m *dtos.RegisterUserRequest) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (user *UserHandler) InsertOne(c echo.Context) error {
	var (
		usr dtos.RegisterUserRequest
		err error
	)

	err = c.Bind(&usr)
	if err != nil {
		return c.JSON(
			http.StatusUnprocessableEntity,
			dtos.NewErrorResponse(
				http.StatusUnprocessableEntity,
				"Filed Cannot Be Empty",
				dtos.GetErrorData(err),
			),
		)
	}

	var ok bool
	if ok, err = isRequestValid(&usr); !ok {
		return c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Bad Request",
				dtos.GetErrorData(err),
			),
		)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := user.UsrUsecase.InsertOne(ctx, &usr)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dtos.NewErrorResponse(
				http.StatusInternalServerError,
				"Cannot Insert Data",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Success",
			result,
		),
	)
}

func (user *UserHandler) FindOne(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := user.UsrUsecase.FindOne(ctx, id)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dtos.NewErrorResponse(
				http.StatusInternalServerError,
				"Cannot Find Data",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Success",
			result,
		),
	)
}

func (user *UserHandler) GetAll(c echo.Context) error {

	type Response struct {
		Total       int64         `json:"total"`
		PerPage     int64         `json:"per_page"`
		CurrentPage int64         `json:"current_page"`
		LastPage    int64         `json:"last_page"`
		From        int64         `json:"from"`
		To          int64         `json:"to"`
		User        []domain.User `json:"users"`
	}

	var (
		res   []domain.User
		count int64
	)

	rp, err := strconv.ParseInt(c.QueryParam("rp"), 10, 64)
	if err != nil {
		rp = 25
	}

	page, err := strconv.ParseInt(c.QueryParam("p"), 10, 64)
	if err != nil {
		page = 1
	}

	filters := bson.D{{Key: "name", Value: primitive.Regex{Pattern: ".*" + c.QueryParam("name") + ".*", Options: "i"}}}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, count, err = user.UsrUsecase.GetAllWithPage(ctx, rp, page, filters, nil)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dtos.NewErrorResponse(
				http.StatusInternalServerError,
				"Cannot Find Data",
				dtos.GetErrorData(err),
			),
		)
	}

	result := Response{
		Total:       count,
		PerPage:     rp,
		CurrentPage: page,
		LastPage:    int64(math.Ceil(float64(count) / float64(rp))),
		From:        page*rp - rp + 1,
		To:          page * rp,
		User:        res,
	}

	return c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Success",
			result,
		),
	)
}

func (user *UserHandler) UpdateOne(c echo.Context) error {
	id := c.Param("id")

	var (
		usr domain.User
		err error
	)

	err = c.Bind(&usr)
	if err != nil {
		return c.JSON(
			http.StatusUnprocessableEntity,
			dtos.NewErrorResponse(
				http.StatusUnprocessableEntity,
				"Filed Cannot Be Empty",
				dtos.GetErrorData(err),
			),
		)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := user.UsrUsecase.UpdateOne(ctx, &usr, id)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dtos.NewErrorResponse(
				http.StatusInternalServerError,
				"Cannot Update Data",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Success",
			result,
		),
	)
}

func (user *UserHandler) DeleteOne(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := user.UsrUsecase.DeleteOne(ctx, id)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dtos.NewErrorResponse(
				http.StatusInternalServerError,
				"Cannot Delete Data",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		dtos.NewResponseMessage(
			http.StatusOK,
			"Success",
		),
	)
}
