package server

import (
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/rielj/go-chatters/pkg/handler"
	"github.com/rielj/go-chatters/pkg/middleware"
	"github.com/rielj/go-chatters/pkg/web"
)

var secretKey = os.Getenv("JWT_SECRET_KEY")

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	fileRoutes(e)

	handler := handler.NewHandler(&handler.Handler{
		Database:   s.db,
		Auth:       s.auth,
		Repository: s.repository,
	})

	r := e.Group("")
	r.Use(middleware.CSP(),
		middleware.JWT(),
		middleware.TextHTML(),
		echomiddleware.Logger(),
		echomiddleware.RemoveTrailingSlash(),
		echomiddleware.Recover(),
	)
	r.GET("/", handler.GetHomeHandler())
	r.GET("/login", handler.GetLoginHandler())
	r.GET("/register", handler.GetRegisterHandler())
	r.GET("/health", handler.GetHealthHandler())

	api := r.Group("/api")
	api.GET("/users", handler.GetUserHandler())
	api.DELETE("/users", handler.DeleteUserHandler())
	api.POST("/users", handler.PostUserHandler())
	api.POST("/login", handler.PostLoginHandler())
	api.POST("/logout", handler.PostLogoutHandler())
	api.POST("/register", handler.PostRegisterHandler())

	return e
}

func fileRoutes(e *echo.Echo) {
	fileServerJS := http.FileServer(http.FS(web.JSFiles))
	fileServerCSS := http.FileServer(http.FS(web.CSSFiles))
	e.GET("/static/css/*", echo.WrapHandler(fileServerCSS))
	e.GET("/static/js/*", echo.WrapHandler(fileServerJS))
}
