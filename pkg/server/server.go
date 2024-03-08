package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/rielj/go-chatters/pkg/auth"
	"github.com/rielj/go-chatters/pkg/database"
	"github.com/rielj/go-chatters/pkg/repository"
)

type Server struct {
	port int

	db         database.Service
	auth       auth.TokenAuth
	repository repository.Repository
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()
	auth := auth.New()
	repository := repository.New(&db)
	NewServer := &Server{
		port: port,

		db:         db,
		auth:       *auth,
		repository: repository,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
