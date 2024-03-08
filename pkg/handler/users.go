package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/rielj/go-chatters/pkg/database"
)

func (h *Handler) GetUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := h.Repository.User.GetUsers()
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{"error": "internal server error"},
			)
		}
		return c.JSON(http.StatusOK, users)
	}
}

func (h *Handler) DeleteUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Param("username")
		err := h.Repository.User.DeleteUser(username)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{"error": "internal server error"},
			)
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "user deleted"})
	}
}

func (h *Handler) PostUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := database.User{}
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
		}
		_, err := h.Repository.User.GetUserByUsername(user.Username)
		if err == nil {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{"error": "username already exists"},
			)
		}
		_, err = h.Repository.User.CreateUser(user)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{"error": "internal server error"},
			)
		}
		return c.JSON(http.StatusOK, user)
	}
}
