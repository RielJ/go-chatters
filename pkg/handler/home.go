package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/rielj/go-chatters/pkg/web/pages"
)

func (h *Handler) GetHomeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := getAuthUser(c, h.Auth)
		if err != nil {
			fmt.Println(err)
			return c.Redirect(http.StatusUnauthorized, "/login")
		}
		fmt.Println(user)

		return render(pages.Home(user.FirstName, user.LastName), c)
	}
}
