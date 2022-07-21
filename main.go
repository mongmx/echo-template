package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/mongmx/echo-template/infra"
	"github.com/mongmx/echo-template/middleware"
	"github.com/mongmx/echo-template/domain/auth"
	"github.com/mongmx/echo-template/domain/user"
)

func main() {
	infra.LoadEnv()
	db, err := infra.NewPostgres(infra.Cfg.Postgres)
	if err != nil {
		log.Fatal("db error", err)
	}
	ce, err := infra.NewCasbin(infra.Cfg.Casbin)
	if err != nil {
		log.Fatal("casbin error", err)
	}

	redisClient, err := infra.NewRedis(infra.Cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	e := echo.New()
	e.Pre(echoMiddleware.RemoveTrailingSlash())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.RequestID())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.HideBanner = true
	e.Validator = infra.NewCustomValidator()

	//*** Start routes register ***//
	// example domain.RouteRegister(e, db)
	auth.RouteRegister(e, db, ce)
	user.RoutesRegister(e, db, redisClient)
	//*** End routes register ***//

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"app":     "pay9-api",
			"version": "0.0.1",
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
