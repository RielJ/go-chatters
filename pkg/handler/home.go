package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/rielj/go-chatters/pkg/web/pages"
)

type GetHomeHandler struct{}

func NewGetHomeHandler() *GetHomeHandler {
	return &GetHomeHandler{}
}

func (h *GetHomeHandler) Handle(c echo.Context) error {
	return render(pages.Home(), c)
}
