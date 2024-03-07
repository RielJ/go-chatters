package server

import (
	"fmt"
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
		TokenLookup: "cookie:x-auth-token",
		ErrorHandler: func(c echo.Context, err error) error {
			cookie, er := c.Cookie("x-auth-token")
			fmt.Println("cookie", cookie)
			fmt.Println("error,", er)

			fmt.Println("error", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		},
	}))

	params := handler.HandlerParams{
		Database: s.db,
		Auth:     s.auth,
	}

	guarded.GET("/", handler.NewGetHomeHandler(params).Handle)

	unguarded := e.Group("")
	unguarded.GET("/login", handler.NewGetLoginHandler(params).Handle)
	unguarded.GET("/register", handler.NewGetRegisterHandler(params).Handle)

	e.GET("/health", handler.NewGetHealthHandler(params).Handle)

	api := e.Group("/api")
	api.GET("/users", handler.NewGetUsersHandler(params).Handle)
	api.POST("/login", handler.NewPostLoginHandler(params).Handle)
	api.POST("/register", handler.NewPostRegisterHandler(params).Handle)

	return e
}

func fileRoutes(e *echo.Echo) {
	fileServerJS := http.FileServer(http.FS(web.JSFiles))
	fileServerCSS := http.FileServer(http.FS(web.CSSFiles))
	e.GET("/static/css/*", echo.WrapHandler(fileServerCSS))
	e.GET("/static/js/*", echo.WrapHandler(fileServerJS))
}
