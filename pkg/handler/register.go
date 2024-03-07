package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/rielj/go-chatters/pkg/database"
	"github.com/rielj/go-chatters/pkg/repository"
	"github.com/rielj/go-chatters/pkg/web/pages"
)

type GetRegisterHandler struct{}

func NewGetRegisterHandler(params HandlerParams) *GetRegisterHandler {
	return &GetRegisterHandler{}
}

func (h *GetRegisterHandler) Handle(c echo.Context) error {
	return render(pages.Register(), c)
}

type PostRegisterHandler struct {
	db database.Service
	ur repository.UserRepository
}

func NewPostRegisterHandler(params HandlerParams) *PostRegisterHandler {
	return &PostRegisterHandler{
		db: params.Database,
		ur: repository.NewUserRepository(repository.UserRepositoryParams{
			Database: params.Database,
		}),
	}
}

func (h *PostRegisterHandler) Handle(c echo.Context) error {
	var user database.User
	user.Username = c.FormValue("username")
	user.Password = c.FormValue("password")
	user.Email = c.FormValue("email")
	user.FirstName = c.FormValue("firstname")
	user.LastName = c.FormValue("lastname")

	if user.Password != c.FormValue("confirm-password") {
		return c.HTML(400, "<p>Passwords do not match</p>")
	}

	_, err := h.ur.CreateUser(user)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}

	cmp := pages.RegisterSuccess()
	cmp.Render(c.Request().Context(), c.Response().Writer)
	return c.HTML(201, "")
}
