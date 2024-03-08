package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetHealthHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, h.Database.Health())
	}
}
