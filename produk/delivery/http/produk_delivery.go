package http

import (
	"context"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"warunk-bem/cloudinary/usecase"
	"warunk-bem/domain"
	"warunk-bem/dtos"
	"warunk-bem/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProdukHandler struct {
	ProdukUsecase domain.ProdukUsecase
}

func NewProdukHandler(router *gin.RouterGroup, protectedAdmin *gin.RouterGroup, pu domain.ProdukUsecase) {
	handler := &ProdukHandler{
		ProdukUsecase: pu,
	}

	api := router.Group("/produk")
	protectedAdmin = protectedAdmin.Group("/produk")

	api.GET("", handler.GetAllWithPage)
	api.GET("/:id", handler.FindOne)
	protectedAdmin.POST("", handler.InsertOne)
	protectedAdmin.PUT("/:id", handler.UpdateOne)
	protectedAdmin.DELETE("/:id", handler.DeleteOne)
}

func isRequestValid(m *dtos.InsertProdukRequest) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (cp *ProdukHandler) GetAllWithPage(c *gin.Context) {
	var (
		res   []*dtos.ProdukDetailResponse
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

	res, count, err = cp.ProdukUsecase.GetAllWithPage(ctx, rp, page, filters, nil)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Cannot Get Produk",
				err.Error(),
			),
		)
		return
	}

	result := dtos.GetAllProdukResponse{
		Total:       count,
		PerPage:     rp,
		CurrentPage: page,
		LastPage:    int64(math.Ceil(float64(count) / float64(rp))),
		From:        (page * rp) - rp + 1,
		To:          page * rp,
		Produk:      res,
	}

	c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Success Get Produk",
			result,
		),
	)
}

func (cp *ProdukHandler) InsertOne(c *gin.Context) {
	var (
		req dtos.InsertProdukRequest
		err error
	)
	name := c.PostForm("name")
	detail := c.PostForm("detail")
	price, _ := strconv.Atoi(c.PostForm("price"))
	stock, _ := strconv.Atoi(c.PostForm("stock"))
	category := c.PostForm("category")

	req = dtos.InsertProdukRequest{
		Slug:     name,
		Name:     name,
		Detail:   detail,
		Price:    int64(price),
		Stock:    int64(stock),
		Category: category,
	}

	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dtos.NewErrorResponse(
				http.StatusUnprocessableEntity,
				"Failed to bind Produk",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	var ok bool
	if ok, err = isRequestValid(&req); !ok {
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

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Cannot Upload Image",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to open file",
				dtos.GetErrorData(err),
			),
		)
		return
	}
	defer src.Close()

	re := regexp.MustCompile(`.png|.jpeg|.jpg`)
	if !re.MatchString(file.Filename) {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewResponseMessage(
				http.StatusBadRequest,
				"The provided file format is not allowed. Please upload a JPEG or PNG image",
			),
		)
		return
	}

	uploadUrl, err := usecase.NewMediaUpload().FileUpload(domain.File{File: src})
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			dtos.NewErrorResponse(
				http.StatusInternalServerError,
				"Error uploading photo",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	req.Image = file

	result, err := cp.ProdukUsecase.InsertOne(ctx, &req, uploadUrl)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Cannot Insert Produk",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Success Insert Produk",
			result,
		),
	)
}

func (cp *ProdukHandler) FindOne(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := cp.ProdukUsecase.FindOne(ctx, id)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Cannot find Produk",
				err.Error(),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Success find Produk",
			result,
		),
	)
}

func (cp *ProdukHandler) UpdateOne(c *gin.Context) {
	var req dtos.ProdukUpdateRequest

	id := c.Param("id")

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dtos.NewErrorResponse(
				http.StatusUnprocessableEntity,
				"Filed Cannot Be Empty",
				err.Error(),
			),
		)
		return
	}

	result, err := cp.ProdukUsecase.UpdateOne(ctx, &req, id)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Cannot Update Produk",
				err.Error(),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Success Update Produk",
			result,
		),
	)
}

func (cp *ProdukHandler) DeleteOne(c *gin.Context) {
	idAdmin, err := middlewares.IsAdmin(c)
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

	var req dtos.DeleteProdukRequest

	id := c.Param("id")

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Cannot Delete Produk",
				err.Error(),
			),
		)
		return
	}

	_, err = cp.ProdukUsecase.DeleteOne(ctx, id, idAdmin, req)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Cannot Delete Produk",
				err.Error(),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dtos.NewResponseMessage(
			http.StatusOK,
			"Success Delete Produk",
		),
	)
}
