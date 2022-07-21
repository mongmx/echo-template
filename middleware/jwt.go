package middleware

import (
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type JwtCustomClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func JWT() echo.MiddlewareFunc {
	return echoMiddleware.JWTWithConfig(echoMiddleware.JWTConfig{
		SigningKey: []byte(os.Getenv("JWT_KEY")),
		Claims:     &JwtCustomClaims{},
	})
}

// func sessionJwt(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		sess, err := session.Get("thruster.sess", c)
// 		if err != nil {
// 			return next(c)
// 		}
// 		tokenString, ok := sess.Values["jwt"].(string)
// 		if !ok {
// 			return next(c)
// 		}
// 		token, err := jwt.ParseWithClaims(
// 			tokenString,
// 			&handlers.JwtCustomClaims{},
// 			func(token *jwt.Token) (interface{}, error) {
// 				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 					return nil, errors.New("unexpected signing method")
// 				}
// 				return []byte(os.Getenv("JWT_KEY")), nil
// 			},
// 		)
// 		if err != nil {
// 			return next(c)
// 		}
// 		claims, ok := token.Claims.(*handlers.JwtCustomClaims)
// 		if !(ok && token.Valid) {
// 			return next(c)
// 		}
// 		c.Set("currentUser", claims)
// 		return next(c)
// 	}
// }
