package auth

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/rielj/go-chatters/pkg/database"
)

type Auth interface {
	GenerateToken(user database.User) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}
