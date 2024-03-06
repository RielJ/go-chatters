package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/rielj/go-chatters/pkg/database"
)

type GetHealthHandler struct {
	db database.Service
}

type GetHealthHandlerParams struct {
	Database database.Service
}

func NewGetHealthHandler(params GetHealthHandlerParams) *GetHealthHandler {
	return &GetHealthHandler{
		db: params.Database,
	}
}

func (h *GetHealthHandler) Handle(c echo.Context) error {
	return c.JSON(http.StatusOK, h.db.Health())
}
