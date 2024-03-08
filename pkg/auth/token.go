package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"

	"github.com/rielj/go-chatters/pkg/database"
)

var secretKey = os.Getenv("JWT_SECRET_KEY")

type TokenAuth struct{}

func NewTokenAuth() *TokenAuth {
	return &TokenAuth{}
}

type CustomClaims struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	jwt.RegisteredClaims
}

func (t *TokenAuth) GenerateToken(user database.User) (string, error) {
	payload := map[string]interface{}{
		"username":  user.Username,
		"email":     user.Email,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims(payload),
	)

	token, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (t *TokenAuth) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
