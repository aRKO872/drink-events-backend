package service

import (
	"context"
	"fmt"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/client"
	"github.com/drink-events-backend/pkg/business_logic/rao"
)

func ResendOTPForVerification(
	ctx context.Context,
	user *models.Users,
	input *models.ResendOTP,
) (bool, error) {
	if input.Event == "verify_email" && user.Email == "" {
		return false, fmt.Errorf("please provide email for email verification")
	}

	if input.Event == "verify_phone" && user.Phone == "" {
		return false, fmt.Errorf("please provide phone num for phone verification")
	}

	// Delete existing OTP for mail or phone
	// Send new OTP to mail and Phone

	otpRedisHelper, redisErr := rao.GetRedisAccessOperator()

	if redisErr != nil {
		return false, redisErr
	}

	var removeOTPErr error
	if input.Event == "verify_email" {
		removeOTPErr = otpRedisHelper.RemoveOTP(ctx, user.Email, input.Event)
		go client.GenerateEmailOTPAndSendEmail(user, otpRedisHelper)
	}

	if input.Event == "verify_phone" {
		removeOTPErr = otpRedisHelper.RemoveOTP(ctx, user.Phone, input.Event)

		// TODO - make a goroutine for sending OTP sms

	}

	if removeOTPErr != nil {
		return false, removeOTPErr
	}

	return true, nil
}