package service

import (
	"context"
	"fmt"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/client"
	"github.com/drink-events-backend/pkg/business_logic/rao"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
)

// Function to Send OTP to provided email successfully
func VerifyEmail(
	ctx context.Context, 
	user *models.Users,
) (bool, error) {
	email := user.Email

	if !pkg_helpers.IsValidEmail(email) {
		return false, fmt.Errorf("email invalid")
	}

	// check if an OTP has already been sent for email verification in 5 minute time window? if yes, return false
	otpRedisHelper, redisErr := rao.GetRedisAccessOperator()

	if redisErr != nil {
		return false, redisErr
	}

	otpExists, _, _ := otpRedisHelper.GetOTP(
		ctx,
		"email",
		email,
		"verify_email",
	)

	if otpExists {
		return false, fmt.Errorf("otp already exists in 5 minute time window")
	}

	// Do below in one Goroutine
	// Generate a new OTP and save it as key `"email"_"event-name"`
	// Send OTP to user's email.
	go client.GenerateEmailOTPAndSendEmail(
		user,
		otpRedisHelper,
	)

	return true, nil
}
