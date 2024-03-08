package handler

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"github.com/rielj/go-chatters/pkg/auth"
	"github.com/rielj/go-chatters/pkg/database"
	"github.com/rielj/go-chatters/pkg/repository"
)

type Handler struct {
	Database   database.Service
	Auth       auth.TokenAuth
	Repository repository.Repository
}

func NewHandler(params *Handler) *Handler {
	return &Handler{
		Database:   params.Database,
		Auth:       params.Auth,
		Repository: params.Repository,
	}
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
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)
}

func setHeader(c echo.Context, name, value string) {
	c.Response().Header().Set(name, value)
}

func setJWTToken(c echo.Context, token string) {
	setCookie(c, "x-auth-token", token, time.Now().Add(24*time.Hour))
	setHeader(c, "x-auth-token", token)
}

func clearJWTToken(c echo.Context) {
	setCookie(c, "x-auth-token", "", time.Now().Add(-24*time.Hour))
	setHeader(c, "x-auth-token", "")
}

func getJWTToken(c echo.Context) (string, error) {
	token, err := c.Cookie("x-auth-token")
	if err == nil {
		return token.Value, nil
	}
	headerToken := c.Request().Header.Get("x-auth-token")
	if headerToken != "" {
		return headerToken, nil
	}
	return "", err
}

func getAuthUser(c echo.Context, auth auth.TokenAuth) (*auth.CustomClaims, error) {
	token, err := c.Cookie("x-auth-token")
	if err != nil {
		return nil, err
	}
	return auth.ValidateToken(token.Value)
}
