package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	_warunkHttp "warunk-bem/warunk/delivery/http"
	_warunkRepo "warunk-bem/warunk/repository"
	_warunktUsecase "warunk-bem/warunk/usecase"
	_wishlistHttp "warunk-bem/wishlist/delivery/http"
	_wishlistRepo "warunk-bem/wishlist/repository"
	_wishlistUsecase "warunk-bem/wishlist/usecase"

	docs "warunk-bem/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
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

// 54.179.176.114:8080/api/v1/swagger/index.html
// 54.179.176.114:8080

// @host      54.179.176.114:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.MaxMultipartMemory = 8 << 20
	middlewares := middlewares.InitMiddleware()
	middlewares.Log()
	r.Use(cors.Default())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	cv := &helpers.CustomValidator{Validators: validator.New()}
	r.Use(func(c *gin.Context) {
		c.Set("validator", cv)
		c.Next()
	})

	CONTEXT_TIMEOUT, err := helpers.GetEnvInt("CONTEXT_TIMEOUT")
	if err != nil {
		log.Fatal(err)
	}
	redisclient := author.InitRedisClient()

	timeoutContext := time.Duration(CONTEXT_TIMEOUT) * time.Second
	database := author.App.Mongo.Database(os.Getenv("MONGODB_NAME"))
	userAmountRepo := _userAmountRepo.NewUserAmountRepository(database)

	userRepo := _userRepo.NewUserRepository(database)
	usrUsecase := _userUcase.NewUserUsecase(userRepo, userAmountRepo, redisclient, timeoutContext)

	// Main Routes API
	api := r.Group("/api/v1")
	protected := r.Group("/api/v1")
	protectedAdmin := r.Group("/api/v1")
	protectedAdmin.Use(middlewares.JwtAuthAdminMiddleware())
	protected.Use(middlewares.JwtAuthMiddleware())

	_userHttp.NewUserHandler(api, protected, protectedAdmin, usrUsecase)

	loginUsecase := _authUsecase.NewAuthUsecase(userRepo, redisclient, timeoutContext)
	_authHttp.NewAuthHandler(api, protected, loginUsecase)

	ProdukRepository := _produkRepo.NewProdukRepository(database)
	ProdukUsecase := _produkUsecase.NewProdukUsecase(ProdukRepository, userRepo, redisclient, timeoutContext)
	_produkHttp.NewProdukHandler(api, protectedAdmin, ProdukUsecase)

	KeranjangRepository := _keranjangRepo.NewKeranjangRepository(database)
	KeranjangUsecase := _keranjangUcase.NewKeranjangUsecase(KeranjangRepository, ProdukRepository, userRepo, redisclient, timeoutContext)
	_keranjangHttp.NewKeranjangHandler(protected, protectedAdmin, KeranjangUsecase, ProdukUsecase)

	FavoriteRepository := _favoriteRepo.NewFavoriteRepository(database)
	FavoriteUsecase := _favoriteUsecase.NewFavoriteUsecase(FavoriteRepository, ProdukRepository, userRepo, redisclient, timeoutContext)
	_favoriteHttp.NewFavoriteHandler(protected, protectedAdmin, FavoriteUsecase, ProdukUsecase)

	WishlistRepository := _wishlistRepo.NewWishlistRepository(database)
	WishlistUsecase := _wishlistUsecase.NewWishlistUsecase(WishlistRepository, ProdukRepository, userRepo, redisclient, timeoutContext)
	_wishlistHttp.NewWishlistHandler(protected, protectedAdmin, WishlistUsecase, ProdukUsecase)

	WarunkRepository := _warunkRepo.NewWarunkRepository(database)
	WarunkUsecase := _warunktUsecase.NewWarunkUsecase(WarunkRepository, ProdukRepository, userRepo, redisclient, timeoutContext)
	_warunkHttp.NewWarunkHandler(protectedAdmin, WarunkUsecase, ProdukUsecase)

	TransaksiRepository := _transaksiRepo.NewTransaksiRepository(database)
	TransaksiUsecase := _transaksiUsecase.NewTransaksiUsecase(TransaksiRepository, KeranjangRepository, ProdukRepository, userRepo, userAmountRepo, WarunkRepository, redisclient, timeoutContext)
	_transaksihttp.NewUserHandler(protected, protectedAdmin, TransaksiUsecase)

	DashboardRepository := _dashboardRepo.NewDashboardRepository(database)
	DashboardUsecase := _dashboardUcase.NewDashboardUsecase(DashboardRepository, userRepo, userAmountRepo, ProdukRepository, TransaksiRepository, redisclient, timeoutContext)
	_dashboardHttp.NewDashboardHandler(protected, DashboardUsecase)

	UserAmountUsecase := _userAmountUsecase.NewUserAmountUsecase(userAmountRepo, userRepo, redisclient, timeoutContext)
	_userAmounthttp.NewUserAmountHandler(protectedAdmin, UserAmountUsecase)

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	api.GET("/healthchecker", func(ctx *gin.Context) {
		requestCtx := ctx.Request.Context()

		value, err := redisclient.Get(requestCtx, "test").Result()
		if err == redis.Nil {
			fmt.Println("key: test does not exist")
		} else if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	appPort := fmt.Sprintf(":%s", os.Getenv("SERVER_ADDRESS"))
	log.Fatal(r.Run(appPort))
}
