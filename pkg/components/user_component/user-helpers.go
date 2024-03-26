package pkg_component_user

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserType string	`json:"user_type"`
	UserId string		`json:"user_id"`
	Email string		`json:"email"`
	Phone string		`json:"phone"`
	Name string			`json:"name"`
	jwt.RegisteredClaims
}

func generateToken(key string, expiration time.Duration, user *Users) (string, error) {
	claims := &JWTClaims{
		UserId:   user.Id,
		UserType: user.UserType,
		Name:     user.Name,
		Phone:    user.Phone,
		Email:    user.Email,
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