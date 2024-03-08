package middleware

import (
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var jwtSecretKey = os.Getenv("JWT_SECRET_KEY")

func JWT() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(jwtSecretKey),
		TokenLookup: "cookie:x-auth-token, header:x-auth-token",
		Skipper: func(c echo.Context) bool {
			paths := []string{"/login", "/register", "/health"}
			for _, path := range paths {
				if c.Path() == path {
					return true
				}
			}
			return false
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		},
	})
}
