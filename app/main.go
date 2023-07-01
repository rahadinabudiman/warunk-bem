package main

import (
	"fmt"
	"log"
	"time"
	_authHttp "warunk-bem/auth/delivery/http"
	_authUsecase "warunk-bem/auth/usecase"
	"warunk-bem/author"
	"warunk-bem/helpers"
	"warunk-bem/middlewares"
	_userHttp "warunk-bem/user/delivery/http"
	_userRepo "warunk-bem/user/repository"
	_userUcase "warunk-bem/user/usecase"
	_userAmountRepo "warunk-bem/user_amount/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {
	r := gin.Default()

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

	appPort := fmt.Sprintf(":%s", author.App.Config.GetString("SERVER_ADDRESS"))
	log.Fatal(r.Run(appPort))
}
