package models

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
)

type CommonErrorOutput struct {
	Status   bool   `json:"status"`
	ErrorMsg string `json:"msg"`
}

type VerifyUserEmailInput struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyOTP struct {
	Event string `json:"event" validate:"required,oneof=verify_email verify_phone login_email login_phone"`
	Otp   int    `json:"otp" validate:"required,gte=100000,lte=999999"`
	Email string `json:"email" validate:"email,omitempty"`
	Phone string `json:"phone" validate:"omitempty,min=10,max=10,regex=^[6-9][0-9]*$"`
}

type LogInInput struct {
	LoggedInFrom string `json:"logged_in_from" validate:"required,oneof=email phone"`
	Email        string `json:"email" validate:"omitempty,email"`
	Phone        string `json:"phone" validate:"omitempty,min=10,max=10,regex=^[6-9][0-9]*$"`
}

type ResendOTP struct {
	Event string `json:"event" validate:"required,oneof=verify_email verify_phone login_email login_phone"`
	Email string `json:"email" validate:"email,omitempty"`
	Phone string `json:"phone" validate:"omitempty,min=10,max=10,regex=^[6-9][0-9]*$"`
}

type SignUpInput struct {
	Name  string `json:"name" validate:"min=3"`
	Email string `json:"email" validate:"email"`
	Phone string `json:"phone" validate:"min=10,max=10,regex=^[6-9][0-9]*$"`
	Bio   string `json:"bio" validate:"omitempty"`
}

type SignUpLoginOutput struct {
	Status       bool   `json:"status"`
	ErrorMsg     string `json:"msg"`
	AccessToken  string `json:"access-token"`
	RefreshToken string `json:"refresh-token"`
}

type GetNearbyUsersOutput struct {
	Status      bool           `json:"status"`
	ErrorMsg    string         `json:"msg"`
	NearbyUsers []UsersSelfTag `json:"users"`
}
type JWTClaims struct {
	UserType     string `json:"user_type"`
	UserId       string `json:"user_id"`
	SearchRadius int    `json:"search_radius"`
	jwt.RegisteredClaims
}

type Users struct {
	Id               string  `json:"id" db:"id"`
	Name             string  `json:"name" db:"name"`
	UserType         string  `json:"user_type" db:"user_type"`
	Email            string  `json:"email" db:"email"`
	Phone            string  `json:"phone" db:"phone"`
	ProfilePicture   string  `json:"profile_picture" db:"profile_picture"`
	Bio              string  `json:"bio" db:"bio"`
	CreatedAt        string  `json:"created_at" db:"created_at"`
	UpdatedAt        string  `json:"updated_at" db:"updated_at"`
	Latitude         float64 `json:"latitude" db:"latitude"`
	Longitude        float64 `json:"longitude" db:"longitude"`
	EmailLastChanged string  `json:"email_last_changed" db:"email_last_changed"`
	PhoneLastChanged string  `json:"phone_last_changed" db:"phone_last_changed"`
	SearchRadius     int     `json:"search_radius" db:"search_radius"`
	IsActive         bool    `json:"is_active" db:"is_active"`
}

type UsersSelfTag struct {
	Users
	Self bool `json:"self"`
}

func (u Users) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *Users) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}

type OTP struct {
	OtpNumber int    `json:"number"`
	Event     string `json:"event"`
	Type      string `json:"type"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

func (otp OTP) MarshalBinary() ([]byte, error) {
	return json.Marshal(otp)
}

func (otp *OTP) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &otp)
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

type LocationSetEvent struct {
	Latitude  float64 `json:"latitude" binding:"required,numeric"`
	Longitude float64 `json:"longitude" binding:"required,numeric"`
}

type UserWithLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UserId    string  `json:"user_id"`
}

type TokenizedUserDetails struct {
	UserType     string `json:"user_type"`
	UserId       string `json:"user_id"`
	SearchRadius int    `json:"search_radius"`
}
