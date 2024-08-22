package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type CommonErrorOutput struct {
	Status   bool   `json:"status"`
	ErrorMsg string `json:"msg"`
}

type SignUpLoginOutput struct {
	Status       bool   `json:"status"`
	ErrorMsg     string `json:"msg"`
	AccessToken  string `json:"access-token,omitempty"`
	RefreshToken string `json:"refresh-token,omitempty"`
}

type GetNearbyUsersOutput struct {
	Status      bool           `json:"status"`
	ErrorMsg    string         `json:"msg"`
	NearbyUsers []UsersSelfTag `json:"users,omitempty"`
}

type JWTClaims struct {
	UserType string `json:"user_type"`
	UserId   string `json:"user_id"`
	jwt.RegisteredClaims
}

type ProjectConfig struct {
	DATABASE_URL               string
	PGADMIN_DEFAULT_EMAIL      string
	PGADMIN_DEFAULT_PASSWORD   string
	REDIS_PASSWORD             string
	POSTGRES_USER              string
	POSTGRES_PASSWORD          string
	POSTGRES_DB                string
	SMTP_HOST                  string
	SMTP_PASSWORD              string
	SMTP_PORT                  int
	SMTP_USER                  string
	ACCESS_TOKEN_EXPIRY        int
	REFRESH_TOKEN_EXPIRY       int
	JWT_SECRET_KEY             string
	FE_SOURCE                  string
	FETCH_NEARBY_ACTIVE_PERIOD int
}

type UserWithLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UserId    string  `json:"user_id"`
}

type TokenizedUserDetails struct {
	UserType string `json:"user_type"`
	UserId   string `json:"user_id"`
}

type SavedImageResponse struct {
	Status   bool   `json:"status"`
	ErrorMsg string `json:"msg"`
	FileId   string `json:"file_id,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type SetUserSearchRadiusResponse struct {
	Status   bool   `json:"status"`
	ErrorMsg string `json:"msg"`
	UserData Users  `json:"user_data,omitempty"`
}

type DeletePairRequestResponse struct {
	Status   bool   `json:"status"`
	ErrorMsg string `json:"msg"`
}

type AcceptPairRequestResponse struct {
	Status   bool   `json:"status"`
	ErrorMsg string `json:"msg"`
	FriendRecord Friends `json:"friends_record,omitempty"`
}