package items

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	jwt "github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "encryption123"

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func GenerateToken(userID string) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "http-api",
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
		UserID: userID,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecret))
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}
