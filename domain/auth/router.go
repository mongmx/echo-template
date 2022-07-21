package auth

import (
	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"github.com/mongmx/echo-template/middleware"
	"gorm.io/gorm"
)

func RouteRegister(e *echo.Echo, db *gorm.DB, ce *casbin.Enforcer) {
	r := NewRepo(db)
	h := NewHandler(r)

	g := e.Group("/auth")
	{
		g.POST("/register", h.Register)
		g.POST("/login", h.Login)
		

		g.Use(middleware.JWT())
		// g.Use(middleware.Auth(db, ce))

		g.GET("/profile", h.GetProfile)
		// g.POST("/merchant/:id", h.SelectMerchant)
	}
}
