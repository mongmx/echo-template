package middleware

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func Auth(rd *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid token")
			}
			uid, err := rd.Get(c.Request().Context(), claims.Subject).Result()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			if uid != claims.ID {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			c.Set("user_id", uid)
			return next(c)
		}
	}
}

// func Permission() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			userID, ok := c.Get("user_id").(string)
// 			if !ok || userID == "" {
// 				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
// 			}
//			read permission by userID
// 			return next(c)
// 		}
// 	}
// }
