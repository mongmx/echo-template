package middleware

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Logger() echo.MiddlewareFunc {
	return echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format: `{"level":"info","time":"${time_rfc3339}",` +
			`"method":"${method}","uri":"${uri}","status":${status},` +
			`"error":"${error}","latency":"${latency_human}"}` + "\n",
	})
}
