package models

import "encoding/json"

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