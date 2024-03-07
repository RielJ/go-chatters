package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/rielj/go-chatters/pkg/database"
	"github.com/rielj/go-chatters/pkg/repository"
)

type GetUsersHandler struct {
	Database database.Service
	ur       repository.UserRepository
}

func NewGetUsersHandler(params HandlerParams) *GetUsersHandler {
	return &GetUsersHandler{
		Database: params.Database,
		ur: repository.NewUserRepository(repository.UserRepositoryParams{
			Database: params.Database,
		}),
	}
}

func (h *GetUsersHandler) Handle(c echo.Context) error {
	users, err := h.ur.GetUsers()
	if err != nil {
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, users)
}

type DeleteUsersHandler struct {
	Database database.Service
	ur       repository.UserRepository
}

func NewDeleteUsersHandler(params HandlerParams) *DeleteUsersHandler {
	return &DeleteUsersHandler{
		Database: params.Database,
		ur: repository.NewUserRepository(repository.UserRepositoryParams{
			Database: params.Database,
		}),
	}
}

func (h *DeleteUsersHandler) Handle(c echo.Context) error {
	username := c.Param("username")
	err := h.ur.DeleteUser(username)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, map[string]string{"message": "user deleted"})
}

type PostUsersHandler struct {
	Database database.Service
	ur       repository.UserRepository
}

func NewPostUsersHandler(params HandlerParams) *PostUsersHandler {
	return &PostUsersHandler{
		Database: params.Database,
		ur: repository.NewUserRepository(repository.UserRepositoryParams{
			Database: params.Database,
		}),
	}
}

func (h *PostUsersHandler) Handle(c echo.Context) error {
	user := database.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(400, map[string]string{"error": "bad request"})
	}
	_, err := h.ur.GetUserByUsername(user.Username)
	if err == nil {
		return c.JSON(400, map[string]string{"error": "username already exists"})
	}
	_, err = h.ur.CreateUser(user)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, user)
}
