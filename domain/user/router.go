package user

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/mongmx/echo-template/middleware"
	"gorm.io/gorm"
)

// Router in auth domain
func RoutesRegister(e *echo.Echo, db *gorm.DB, rd *redis.Client) {
	s := NewStore(db)
	h := NewHandler(s, db, rd)
	g := e.Group("/api")
	{
		g.POST("/user/login", h.SignIn)
		g.POST("/user/register", h.SignUp)
		g.POST("/user/refresh-token", h.RefreshToken)

		// TODO: Email verify
		// e.POST("/api/auth/email-verify", h.EmailVerify)
		// TODO: OTP login
		// e.POST("/api/users/otp/request", h.OTPRequest)
		// e.POST("/api/users/otp/verify", h.OTPVerify)

		g.Use(middleware.JWT())
		g.Use(middleware.Auth(rd))

		g.GET("/profile", h.CurrentUser)
		g.GET("/current-user", h.CurrentUser)
		g.POST("/update-password", h.UpdatePassword)
		g.GET("/user/get-history", h.GetHistory)
		g.DELETE("/signout", h.SignOut)
		g.GET("/profiles/:username", h.GetProfile)                // RealWorld: Get a profiles
		g.POST("/profiles/:username/follow", h.FollowProfile)     // RealWorld: Follow a user
		g.DELETE("/profiles/:username/follow", h.UnfollowProfile) // RealWorld: Unfollow a user
		g.PUT("", h.UpdatePassword)                               // TODO: RealWorld:
	}

	e.GET("/api/user", h.CurrentUser) // RealWorld:
}
