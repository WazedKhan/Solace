package jwt_token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"

var (
	ErrTokenCreationFailed = errors.New("failed to generate token")
	ErrInvalidToken        = errors.New("invalid token")
)

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

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(UserIDKey).(string)
	return id, ok
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

func (g *Generator) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return g.secret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
