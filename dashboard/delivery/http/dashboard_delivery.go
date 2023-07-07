package http

import (
	"context"
	"net/http"
	"strconv"
	"warunk-bem/domain"
	"warunk-bem/dtos"
	"warunk-bem/middlewares"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	DashboardUsecase domain.DashboardUsecase
}

func NewDashboardHandler(router *gin.RouterGroup, du domain.DashboardUsecase) {
	handler := &DashboardHandler{
		DashboardUsecase: du,
	}

	api := router.Group("/dashboard")

	api.GET("", handler.GetDashboardData)
}

func (dh *DashboardHandler) GetDashboardData(c *gin.Context) {
	userID, err := middlewares.IsUser(c)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Cannot get Admin",
				err.Error(),
			),
		)
		return
	}

	rp, _ := strconv.ParseInt(c.Query("rp"), 10, 64)
	p, _ := strconv.ParseInt(c.Query("p"), 10, 64)

	filter := make(map[string]interface{})
	// Set filter based on query parameters
	// Example: filter["field"] = value

	// Set sorting options
	// setsort := make(map[string]interface{})
	// Example: setsort["field"] = 1 // 1 for ascending, -1 for descending

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	data, err := dh.DashboardUsecase.GetDashboardData(ctx, userID, rp, p, filter, nil)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get dashboard data",
				err.Error(),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Success",
			data,
		),
	)
}
