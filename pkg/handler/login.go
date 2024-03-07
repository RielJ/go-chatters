package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/rielj/go-chatters/pkg/auth"
	"github.com/rielj/go-chatters/pkg/database"
	"github.com/rielj/go-chatters/pkg/repository"
	"github.com/rielj/go-chatters/pkg/tools"
	"github.com/rielj/go-chatters/pkg/web/pages"
)

type GetLoginHandler struct {
	Auth auth.Auth
}

func NewGetLoginHandler(params HandlerParams) *GetLoginHandler {
	return &GetLoginHandler{
		Auth: params.Auth,
	}
}

func (h *GetLoginHandler) Handle(c echo.Context) error {
	token, err := c.Cookie("x-auth-token")
	if err != nil {
		return render(pages.Login(), c)
	}
	fmt.Println("token", token)
	jwtClaims, err := auth.NewTokenAuth().
		ValidateToken(strings.TrimLeft(token.String(), "x-auth-token="))
	if err != nil {
		return c.Redirect(http.StatusFound, "/")
	}
	fmt.Println("jwt", jwtClaims)

	return render(pages.Login(), c)
}

type PostLoginHandler struct {
	Database database.Service
	Auth     auth.Auth
	ur       repository.UserRepository
}

func NewPostLoginHandler(params HandlerParams) *PostLoginHandler {
	return &PostLoginHandler{
		Database: params.Database,
		Auth:     params.Auth,
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
		fmt.Println("error generating token", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}

	cookie := new(http.Cookie)
	cookie.Name = "x-auth-token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(7 * 24 * time.Hour)
	cookie.Path = "/"
	c.SetCookie(cookie)

	c.Response().Header().Set("HX-Redirect", "/")
	return c.HTML(200, "Logged in successfully!")
}
