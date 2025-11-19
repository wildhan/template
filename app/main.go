package main

import (
	"fmt"
	"net/http"
	"os"

	"template/config/database"
	"template/lib/log"
	"template/lib/token"

	userHandler "template/package/user/handler"
	userRepository "template/package/user/repository"
	userUsecase "template/package/user/usecase"

	authHandler "template/package/auth/handler"
	authRepository "template/package/auth/repository"
	authUsacase "template/package/auth/usecase"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Error(fmt.Sprintf("Failed load .env: %v\n", err.Error()))
		os.Exit(2)
	}

	dbConn := database.CreateConnection()
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{http.MethodGet},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Template")
	})

	tokenMaker, err := token.NewPasetoMaker(os.Getenv("SECRET_KEY_64_BYTES"))
	if err != nil {
		log.Error(fmt.Sprintf("Failed create token maker: %v\n", err.Error()))
		os.Exit(2)
	}

	userRepo := userRepository.NewUserRepo(dbConn)
	userUC := userUsecase.NewUserUsecase(userRepo)
	userHandler.NewUserHandler(userUC).Mount(e.Group("/user"))

	authRepo := authRepository.NewAuthRepo()
	authUC := authUsacase.NewAuthUsecase(dbConn, tokenMaker, authRepo)
	authHandler.NewAuthHandler(authUC).Mount(e.Group("/auth"))

	if err := e.Start(":" + os.Getenv("PORT")); err != nil {
		log.Error(fmt.Sprintf("Failed start echo: %v", err.Error()))
	}
}
