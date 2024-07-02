package service

import (
	"context"
	"fmt"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/rao"
)

// Function to Send OTP to provided email successfully
func VerifyOTP(
	ctx context.Context,
	otp int, 
	event string, 
	user *models.Users,
) (bool, error) {
	switch event {
	case "verify_email":
		return verifyOTPForEmail(ctx, otp, event, user)
	case "verify_phone":
		// TODO when otp for phone is set up
		return true, nil
	default:
		return false, fmt.Errorf("please provide proper event")
	}
}

func verifyOTPForEmail(
	ctx context.Context,
	otp int, 
	event string,
	user *models.Users,
) (bool, error) {
	email := user.Email

	roo, redisFetchErr := rao.GetRedisAccessOperator()

	if redisFetchErr != nil {
		return false, redisFetchErr
	}

	_, fetchedOTP, otpFetchErr := roo.GetOTP(ctx, "email", email, event)

	if otpFetchErr != nil {
		return false, otpFetchErr
	}

	if fetchedOTP.OtpNumber != otp {
		return false, fmt.Errorf("wrong otp provided")
	}

	if removeOTPErr := roo.RemoveOTP(ctx, user.Email, event); removeOTPErr != nil {
		return false, fmt.Errorf("error removing otp : %s", removeOTPErr.Error())
	}

	return true, nil
}