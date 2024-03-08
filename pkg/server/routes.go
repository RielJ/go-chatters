package server

import (
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/rielj/go-chatters/pkg/handler"
	"github.com/rielj/go-chatters/pkg/web"
)

var secretKey = os.Getenv("JWT_SECRET_KEY")

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	fileRoutes(e)

	e.Use(middleware.Logger())
	e.Use(middleware.RemoveTrailingSlash())

	guarded := e.Group("")
	guarded.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(secretKey),
		TokenLookup: "cookie:x-auth-token, header:x-auth-token",
		ErrorHandler: func(c echo.Context, err error) error {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		},
	}))

	handler := handler.NewHandler(&handler.Handler{
		Database:   s.db,
		Auth:       s.auth,
		Repository: s.repository,
	})

	guarded.GET("/", handler.GetHomeHandler())

	unguarded := e.Group("")
	unguarded.GET("/login", handler.GetLoginHandler())
	unguarded.GET("/register", handler.GetRegisterHandler())

	e.GET("/health", handler.GetHealthHandler())

	api := e.Group("/api")
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
