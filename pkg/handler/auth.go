package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/rielj/go-chatters/pkg/database"
	"github.com/rielj/go-chatters/pkg/tools"
	"github.com/rielj/go-chatters/pkg/web/pages"
)

func (h *Handler) GetLoginHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := c.Cookie("x-auth-token")
		if err != nil {
			return render(pages.Login(), c)
		}
		fmt.Println("token", token)
		jwtClaims, err := h.Auth.ValidateToken(
			strings.TrimLeft(token.String(), "x-auth-token="),
		)
		if err != nil {
			return c.Redirect(http.StatusFound, "/")
		}
		fmt.Println("jwt", jwtClaims)

		return render(pages.Login(), c)
	}
}

func (h *Handler) PostLoginHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		user, err := h.Repository.User.GetUserByUsername(username)
		if err != nil {
			fmt.Println("user not found", err)
			pages.LoginError().Render(c.Request().Context(), c.Response().Writer)
			return c.HTML(http.StatusUnauthorized, "")
		}

		err = tools.CheckPasswordHash(user.Password, password)
		if err != nil {
			fmt.Println("password incorrect", err)
			pages.LoginError().Render(c.Request().Context(), c.Response().Writer)
			return c.HTML(http.StatusUnauthorized, "")
		}

		token, err := h.Auth.GenerateToken(*user)
		if err != nil {
			fmt.Println("error generating token", err)
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{"error": "internal server error"},
			)
		}

		cookie := new(http.Cookie)
		cookie.Name = "x-auth-token"
		cookie.Value = token
		cookie.Expires = time.Now().Add(7 * 24 * time.Hour)
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.Secure = true
		cookie.SameSite = http.SameSiteLaxMode
		c.SetCookie(cookie)

		c.Response().Header().Set("HX-Redirect", "/")
		return c.HTML(http.StatusOK, "Logged in successfully!")
	}
}

func (h *Handler) GetRegisterHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return render(pages.Register(), c)
	}
}

func (h *Handler) PostRegisterHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var user database.User
		user.Username = c.FormValue("username")
		user.Password = c.FormValue("password")
		user.Email = c.FormValue("email")
		user.FirstName = c.FormValue("firstname")
		user.LastName = c.FormValue("lastname")

		err := h.Repository.User.ValidateUserFields(user)
		if err != nil {
			pages.RegisterError(err.Error()).Render(c.Request().Context(), c.Response().Writer)
			return c.HTML(http.StatusBadRequest, "")
		}

		if user.Password != c.FormValue("confirm-password") {
			pages.RegisterError("Passwords do not match").
				Render(c.Request().Context(), c.Response().Writer)
			return c.HTML(http.StatusBadRequest, "")
		}

		_, err = h.Repository.User.GetUserByUsername(user.Username)
		if err == nil {
			pages.RegisterError("Username already exists").
				Render(c.Request().Context(), c.Response().Writer)
			return c.HTML(http.StatusBadRequest, "")
		}

		_, err = h.Repository.User.GetUserByEmail(user.Email)
		if err == nil {
			pages.RegisterError("Email already exists").
				Render(c.Request().Context(), c.Response().Writer)
			return c.HTML(http.StatusBadRequest, "")
		}

		_, err = h.Repository.User.CreateUser(user)
		if err != nil {
			fmt.Println("ERROR: ", err)
			return err
		}

		cmp := pages.RegisterSuccess()
		cmp.Render(c.Request().Context(), c.Response().Writer)
		return c.HTML(http.StatusCreated, "")
	}
}

func (h *Handler) PostLogoutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie := new(http.Cookie)
		cookie.Name = "x-auth-token"
		cookie.Value = ""
		cookie.Expires = time.Now()
		cookie.Path = "/"
		c.SetCookie(cookie)

		return c.Redirect(http.StatusFound, "/login")
	}
}
