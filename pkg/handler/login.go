package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/rielj/go-chatters/pkg/auth"
	"github.com/rielj/go-chatters/pkg/database"
	"github.com/rielj/go-chatters/pkg/repository"
	"github.com/rielj/go-chatters/pkg/tools"
	"github.com/rielj/go-chatters/pkg/web/pages"
)

type GetLoginHandler struct{}

func NewGetLoginHandler() *GetLoginHandler {
	return &GetLoginHandler{}
}

func (h *GetLoginHandler) Handle(c echo.Context) error {
	return render(pages.Login(), c)
}

type PostLoginHandler struct {
	Database database.Service
	Auth     auth.Auth
	ur       repository.UserRepository
}

type LoginHandlerParams struct {
	Database      database.Service
	Authenticator auth.Auth
}

func NewPostLoginHandler(params LoginHandlerParams) *PostLoginHandler {
	return &PostLoginHandler{
		Database: params.Database,
		Auth:     params.Authenticator,
		ur: repository.NewUserRepository(
			repository.UserRepositoryParams{
				Database: params.Database,
			},
		),
	}
}

func (h *PostLoginHandler) Handle(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.ur.GetUserByUsername(username)
	if err != nil {
		fmt.Println("user not found", err)
		pages.LoginError().Render(c.Request().Context(), c.Response().Writer)
		return c.HTML(401, "")
	}

	err = tools.CheckPasswordHash(user.Password, password)
	if err != nil {
		fmt.Println("password incorrect", err)
		pages.LoginError().Render(c.Request().Context(), c.Response().Writer)
		return c.HTML(401, "")
	}

	token, err := h.Auth.GenerateToken(*user)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}

	c.SetCookie(&http.Cookie{
		Name:    "access_token",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(7 * 24 * time.Hour),
	})

	c.Response().Header().Set("HX-Redirect", "/")
	return c.HTML(200, "Logged in successfully!")
}
