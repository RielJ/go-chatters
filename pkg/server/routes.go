package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/rielj/go-chatters/pkg/handler"
	"github.com/rielj/go-chatters/pkg/web"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())

	fileRoutes(e)

	e.GET("/", handler.NewGetHomeHandler().Handle)

	e.GET("/login", handler.NewGetLoginHandler().Handle)
	e.GET("/register", handler.NewGetRegisterHandler().Handle)
	e.GET("/health", handler.NewGetHealthHandler(handler.GetHealthHandlerParams{
		Database: s.db,
	}).Handle)

	api := e.Group("/api")
	api.GET("/users", handler.NewGetUsersHandler(handler.UsersHandlerParams{
		Database: s.db,
	}).Handle)
	api.POST("/login", handler.NewPostLoginHandler(handler.LoginHandlerParams{
		Database:      s.db,
		Authenticator: s.auth,
	}).Handle)
	api.POST("/register", handler.NewPostRegisterHandler(handler.PostRegisterHandlerParams{
		Database: s.db,
	}).Handle)

	return e
}

func fileRoutes(e *echo.Echo) {
	fileServerJS := http.FileServer(http.FS(web.JSFiles))
	fileServerCSS := http.FileServer(http.FS(web.CSSFiles))
	e.GET("/static/css/*", echo.WrapHandler(fileServerCSS))
	e.GET("/static/js/*", echo.WrapHandler(fileServerJS))
}
