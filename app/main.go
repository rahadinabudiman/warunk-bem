package main

import (
	"fmt"
	"log"
	"time"
	"warunk-bem/author"
	_userHttp "warunk-bem/user/delivery/http"
	_userHttpMiddlewares "warunk-bem/user/delivery/http/middlewares"
	_userRepo "warunk-bem/user/repository"
	_userUcase "warunk-bem/user/usecase"
	"warunk-bem/user/usecase/helpers"

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
	userRepo := _userRepo.NewUserRepository(database)
	usrUsecase := _userUcase.NewUserUsecase(userRepo, timeoutContext)
	_userHttp.NewUserHandler(e, usrUsecase)

	appPort := fmt.Sprintf(":%s", author.App.Config.GetString("SERVER_ADDRESS"))
	log.Fatal(e.Start(appPort))
}
