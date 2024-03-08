package middleware

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func New(e *echo.Echo) {
	e.Use(
		CSP(),
		JWT(),
		echoMiddleware.Logger(),
		echoMiddleware.RemoveTrailingSlash(),
		echoMiddleware.Recover(),
	)
}
