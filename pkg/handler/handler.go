package handler

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Handle(c echo.Context) echo.HandlerFunc
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
