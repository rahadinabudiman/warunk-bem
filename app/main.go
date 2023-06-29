package main

import (
	"fmt"
	"log"
	"time"
	_authHttp "warunk-bem/auth/delivery/http"
	_authUsecase "warunk-bem/auth/usecase"
	"warunk-bem/author"
	_jwt "warunk-bem/jwt/usecase"
	_userHttp "warunk-bem/user/delivery/http"
	_userHttpMiddlewares "warunk-bem/user/delivery/http/middlewares"
	_userRepo "warunk-bem/user/repository"
	_userUcase "warunk-bem/user/usecase"
	"warunk-bem/user/usecase/helpers"
	_userAmountRepo "warunk-bem/user_amount/repository"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	middlewares := _userHttpMiddlewares.InitMiddleware()
	middlewares.Log(e)
	e.Use(middlewares.AllowCORS)
	e.Use(middlewares.CORS)

	cv := &helpers.CustomValidator{Validators: validator.New()}
	e.Validator = cv

	timeoutContext := time.Duration(author.App.Config.GetInt("CONTEXT_TIMEOUT")) * time.Second
	database := author.App.Mongo.Database(author.App.Config.GetString("MONGODB_NAME"))
	userAmountRepo := _userAmountRepo.NewUserAmountRepository(database)

	userRepo := _userRepo.NewUserRepository(database)
	usrUsecase := _userUcase.NewUserUsecase(userRepo, userAmountRepo, timeoutContext)
	_userHttp.NewUserHandler(e, usrUsecase)

	// Main Routes API
	api := e.Group("/api/v1")

	jwt := _jwt.NewJwtUsecase(userRepo, timeoutContext, author.App.Config)
	userJwt := api.Group("")
	jwt.SetJwtUser(userJwt)
	adminJwt := api.Group("")
	jwt.SetJwtUser(adminJwt)
	generalJwt := api.Group("")
	jwt.SetJwtUser(generalJwt)

	loginUsecase := _authUsecase.NewAuthUsecase(userRepo, timeoutContext, author.App.Config)
	_authHttp.NewAuthHandler(api, generalJwt, loginUsecase, author.App.Config)

	appPort := fmt.Sprintf(":%s", author.App.Config.GetString("SERVER_ADDRESS"))
	log.Fatal(e.Start(appPort))
}
