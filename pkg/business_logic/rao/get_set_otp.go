package rao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/drink-events-backend/models"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
	"github.com/redis/go-redis/v9"
)

func (rao *RedisAccessOperator) GetOTP(
	ctx context.Context,
	otpType string,
	emailOrPhone string,
	eventType string,
) (bool, *models.OTP, error) {

	rdbClient := rao.RDB

	var otp *models.OTP

	if otpType == "phone" && !pkg_helpers.IsValidPhoneNumber(emailOrPhone) {
		return false, nil, fmt.Errorf("not a valid phone number")
	} else if otpType == "email" && !pkg_helpers.IsValidEmail(emailOrPhone) {
		return false, nil, fmt.Errorf("not a valid email")
	} else if otpType != "phone" && otpType != "email" {
		return false, nil, fmt.Errorf("invalid otp type provided")
	}

	// Construct a Key to save OTP data
	key := fmt.Sprintf("%s_%s", emailOrPhone, eventType)

	// Checking if exists user exists 
	otpExists, existErr := rdbClient.Exists(ctx, key).Result()

	if existErr != nil {
		return false, nil, fmt.Errorf("error checking existence of OTP: %s", existErr.Error())
	}

	if otpExists != 1 {
		return false, nil, fmt.Errorf("otp does not exist")
	}

	// user exists and fetching and putting value in User
	val, err := rdbClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil, nil // Key does not exist
		}
		return false, nil, fmt.Errorf("error fetching OTP data: %s", err.Error())
	}

	// Unmarshal OTP data
	if err := json.Unmarshal([]byte(val), &otp); err != nil {
		return false, nil, fmt.Errorf("error unmarshalling OTP data: %s", err.Error())
	}

	return true, otp, nil
}

func (rao *RedisAccessOperator) SetOTP(
	ctx context.Context,
	otpNumber int,
	otpType string,
	emailOrPhone string,
	eventType string,
	noOfMinutes int,
) error {
	if otpType == "phone" && !pkg_helpers.IsValidPhoneNumber(emailOrPhone) {
		return fmt.Errorf("not a valid phone number")
	} else if otpType == "email" && !pkg_helpers.IsValidEmail(emailOrPhone) {
		return fmt.Errorf("not a valid email")
	} else if otpType != "phone" && otpType != "email" {
		return fmt.Errorf("invalid otp type provided")
	}

	rdbClient := rao.RDB

	// Setting No of minutes as expiration time in Redis
	rao.Lock()
	defer rao.Unlock()
	otpObj := models.OTP{
		OtpNumber: otpNumber,
		Event: eventType,
		Type: otpType,
		Email: emailOrPhone,
		Phone: emailOrPhone,
	}
	setErr := rdbClient.Set(
		ctx, 
		fmt.Sprintf("%s_%s", emailOrPhone, eventType), 
		otpObj,
		(time.Duration(noOfMinutes) * time.Minute),
	).Err()

	if setErr != nil {
		return setErr
	}

	return nil
}

func (rao *RedisAccessOperator) RemoveOTP(
	ctx context.Context,
	emailOrPhone string,
	event string,
) error {
	rdbClient := rao.RDB

	key := fmt.Sprintf("%s_%s", emailOrPhone, event)

	rao.Lock()
	defer rao.Unlock()
	if err := rdbClient.Del(
		ctx, 
		key,
	).Err(); err != nil {
		return fmt.Errorf("failed to remove otp")
	}

	return nil
}