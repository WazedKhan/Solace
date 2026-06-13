package jwt_token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrTokenCreationFailed = errors.New("failed to generate token")

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type Generator struct {
	secret []byte
	ttl    time.Duration
}

func NewGenerator(secret string, ttl time.Duration) *Generator {
	return &Generator{secret: []byte(secret), ttl: ttl}
}

func (g *Generator) Generate(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "solace-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(g.secret)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrTokenCreationFailed, err)
	}
	return tokenString, nil
}
