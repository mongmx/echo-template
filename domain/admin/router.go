package admin

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RoutesRegister(e *echo.Echo,db *gorm.DB, rd *redis.Client) {
	s := NewStore(db)
	h := NewHandler(s, rd)

	g := e.Group("/api/admin")
	{
		g.GET("/profiles", h.ListProfiles)
	}
}
