package infrastructure

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager interface {
	Generate(userID string) (string, error)
	Validate(token string) (string, error)
}

type jwtManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret string, ttl time.Duration) JWTManager {
	return &jwtManager{secret: []byte(secret), ttl: ttl}
}

func (j *jwtManager) Generate(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(j.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *jwtManager) Validate(tokenStr string) (string, error) {
	if tokenStr == "" {
		return "", errors.New("empty token")
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims["sub"].(string), nil
}
