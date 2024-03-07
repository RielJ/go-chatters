package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/rielj/go-chatters/pkg/auth"
	"github.com/rielj/go-chatters/pkg/web/pages"
)

type GetHomeHandler struct {
	Auth auth.Auth
}

func NewGetHomeHandler(params HandlerParams) *GetHomeHandler {
	return &GetHomeHandler{
		Auth: params.Auth,
	}
}

func (h *GetHomeHandler) Handle(c echo.Context) error {
	user, err := getAuthUser(c, h.Auth)
	if err != nil {
		return c.Redirect(http.StatusUnauthorized, "/login")
	}

	fmt.Println("user", user)

	return render(pages.Home(), c)
}
