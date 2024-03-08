package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetHealthHandler(params HandlerParams) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, params.Database.Health())
	}
}
