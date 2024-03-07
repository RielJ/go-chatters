package handler

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/rielj/go-chatters/pkg/auth"
	"github.com/rielj/go-chatters/pkg/database"
)

type Handler interface {
	Handle(c echo.Context) echo.HandlerFunc
}

type HandlerParams struct {
	Database database.Service
	Auth     auth.Auth
}

func render(comp templ.Component, c echo.Context) error {
	return echo.WrapHandler(templ.Handler(comp))(c)
}

func setCookie(c echo.Context, name, value string, expiresAt time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Expires = expiresAt
	cookie.Path = "/"
	c.SetCookie(cookie)
}

func getAuthUser(c echo.Context, auth auth.Auth) (jwt.MapClaims, error) {
	token, err := c.Cookie("x-auth-token")
	if err != nil {
		return nil, err
	}
	claims, err := auth.ValidateToken(token.Value)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
