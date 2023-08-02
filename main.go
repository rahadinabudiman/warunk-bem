package main

import (
	"fmt"
	"log"
	"time"
	_authHttp "warunk-bem/auth/delivery/http"
	_authUsecase "warunk-bem/auth/usecase"
	"warunk-bem/author"
	_dashboardHttp "warunk-bem/dashboard/delivery/http"
	_dashboardRepo "warunk-bem/dashboard/repository"
	_dashboardUcase "warunk-bem/dashboard/usecase"
	_favoriteHttp "warunk-bem/favorite/delivery/http"
	_favoriteRepo "warunk-bem/favorite/repository"
	_favoriteUsecase "warunk-bem/favorite/usecase"
	"warunk-bem/helpers"
	_keranjangHttp "warunk-bem/keranjang/delivery/http"
	_keranjangRepo "warunk-bem/keranjang/repository"
	_keranjangUcase "warunk-bem/keranjang/usecase"
	"warunk-bem/middlewares"
	_produkHttp "warunk-bem/produk/delivery/http"
	_produkRepo "warunk-bem/produk/repository"
	_produkUsecase "warunk-bem/produk/usecase"
	_transaksihttp "warunk-bem/transaksi/delivery/http"
	_transaksiRepo "warunk-bem/transaksi/repository"
	_transaksiUsecase "warunk-bem/transaksi/usecase"
	_userHttp "warunk-bem/user/delivery/http"
	_userRepo "warunk-bem/user/repository"
	_userUcase "warunk-bem/user/usecase"
	_userAmounthttp "warunk-bem/user_amount/delivery/http"
	_userAmountRepo "warunk-bem/user_amount/repository"
	_userAmountUsecase "warunk-bem/user_amount/usecase"
	_wishlistHttp "warunk-bem/wishlist/delivery/http"
	_wishlistRepo "warunk-bem/wishlist/repository"
	_wishlistUsecase "warunk-bem/wishlist/usecase"

	docs "warunk-bem/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Warunk BEM Documentation API
// @version         1.0
// @termsOfService  http://swagger.io/terms/

// @contact.name   r4ha
// @contact.url    https://github.com/rahadinabudiman

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// localhost:8080/api/v1/swagger/index.html
// localhost:8080

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.MaxMultipartMemory = 8 << 20
	middlewares := middlewares.InitMiddleware()
	middlewares.Log()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://keyzex.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	cv := &helpers.CustomValidator{Validators: validator.New()}
	r.Use(func(c *gin.Context) {
		c.Set("validator", cv)
		c.Next()
	})

	timeoutContext := time.Duration(author.App.Config.GetInt("CONTEXT_TIMEOUT")) * time.Second
	database := author.App.Mongo.Database(author.App.Config.GetString("MONGODB_NAME"))
	userAmountRepo := _userAmountRepo.NewUserAmountRepository(database)

	userRepo := _userRepo.NewUserRepository(database)
	usrUsecase := _userUcase.NewUserUsecase(userRepo, userAmountRepo, timeoutContext)

	// Main Routes API
	api := r.Group("/api/v1")
	protected := r.Group("/api/v1")
	protectedAdmin := r.Group("/api/v1")
	protectedAdmin.Use(middlewares.JwtAuthAdminMiddleware())
	protected.Use(middlewares.JwtAuthMiddleware())

	_userHttp.NewUserHandler(api, protected, protectedAdmin, usrUsecase)

	loginUsecase := _authUsecase.NewAuthUsecase(userRepo, timeoutContext, author.App.Config)
	_authHttp.NewAuthHandler(api, protected, loginUsecase, author.App.Config)

	ProdukRepository := _produkRepo.NewProdukRepository(database)
	ProdukUsecase := _produkUsecase.NewProdukUsecase(ProdukRepository, userRepo, timeoutContext)
	_produkHttp.NewProdukHandler(api, protectedAdmin, ProdukUsecase)

	KeranjangRepository := _keranjangRepo.NewKeranjangRepository(database)
	KeranjangUsecase := _keranjangUcase.NewKeranjangUsecase(KeranjangRepository, ProdukRepository, userRepo, timeoutContext)
	_keranjangHttp.NewKeranjangHandler(protected, protectedAdmin, KeranjangUsecase, ProdukUsecase)

	TransaksiRepository := _transaksiRepo.NewTransaksiRepository(database)
	TransaksiUsecase := _transaksiUsecase.NewTransaksiUsecase(TransaksiRepository, KeranjangRepository, ProdukRepository, userRepo, userAmountRepo, timeoutContext)
	_transaksihttp.NewUserHandler(protected, protectedAdmin, TransaksiUsecase)

	DashboardRepository := _dashboardRepo.NewDashboardRepository(database)
	DashboardUsecase := _dashboardUcase.NewDashboardUsecase(DashboardRepository, userRepo, userAmountRepo, ProdukRepository, TransaksiRepository, timeoutContext)
	_dashboardHttp.NewDashboardHandler(protected, DashboardUsecase)

	FavoriteRepository := _favoriteRepo.NewFavoriteRepository(database)
	FavoriteUsecase := _favoriteUsecase.NewFavoriteUsecase(FavoriteRepository, ProdukRepository, userRepo, timeoutContext)
	_favoriteHttp.NewFavoriteHandler(protected, protectedAdmin, FavoriteUsecase, ProdukUsecase)

	WishlistRepository := _wishlistRepo.NewWishlistRepository(database)
	WishlistUsecase := _wishlistUsecase.NewWishlistUsecase(WishlistRepository, ProdukRepository, userRepo, timeoutContext)
	_wishlistHttp.NewWishlistHandler(protected, protectedAdmin, WishlistUsecase, ProdukUsecase)

	UserAmountUsecase := _userAmountUsecase.NewUserAmountUsecase(userAmountRepo, userRepo, timeoutContext)
	_userAmounthttp.NewUserAmountHandler(protectedAdmin, UserAmountUsecase)

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	appPort := fmt.Sprintf(":%s", author.App.Config.GetString("SERVER_ADDRESS"))
	log.Fatal(r.Run(appPort))
}
