package pkg_helpers

import (
	"errors"
	"fmt"
	"time"

	"github.com/drink-events-backend/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(key string, expiration time.Duration, user *models.Users) (string, error) {
	claims := &models.JWTClaims{
		UserId:   user.Id,
		UserType: user.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key[:]))
	if err != nil {
		return "", errors.New("error occured while generating token")
	}

	return signedToken, nil
}

func ParseToken(tokenString string, key string) (*models.JWTClaims, error) {
	claims := &models.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		fmt.Println("Error is : ", err)
		return nil, errors.New("error occurred while parsing token")
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	return claims, nil
}

