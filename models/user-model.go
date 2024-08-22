package models

import "encoding/json"

type Users struct {
	Id               string  `json:"id,omitempty" db:"id"`
	Name             string  `json:"name,omitempty" db:"name"`
	UserType         string  `json:"user_type,omitempty" db:"user_type"`
	Email            string  `json:"email,omitempty" db:"email"`
	Phone            string  `json:"phone,omitempty" db:"phone"`
	ProfilePicture   string  `json:"profile_picture,omitempty" db:"profile_picture"`
	Bio              string  `json:"bio,omitempty" db:"bio"`
	CreatedAt        string  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt        string  `json:"updated_at,omitempty" db:"updated_at"`
	Latitude         float64 `json:"latitude,omitempty" db:"latitude"`
	Longitude        float64 `json:"longitude,omitempty" db:"longitude"`
	EmailLastChanged string  `json:"email_last_changed,omitempty" db:"email_last_changed"`
	PhoneLastChanged string  `json:"phone_last_changed,omitempty" db:"phone_last_changed"`
	SearchRadius     int     `json:"search_radius,omitempty" db:"search_radius"`
	IsActive         bool    `json:"is_active,omitempty" db:"is_active"`
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