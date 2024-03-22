package pkg_component

import (
	"context"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	pkg_helpers "github.com/drink-events-backend/pkg/helper"
	redis_operations "github.com/drink-events-backend/pkg/redis-operations"
	"github.com/redis/go-redis/v9"
)

type RedisOTPOperator struct {
	redis_operations.RedisOperator
}

func GetRedisOTPOperator() (*RedisOTPOperator, error) {
	ro, initErr := redis_operations.CreateOrGetRedisOperator()

	if initErr != nil {
		return nil, initErr
	}

	return &RedisOTPOperator{
		RedisOperator: *ro,
	}, nil
}

func (roo *RedisOTPOperator) Get(
	otpType string,
	emailOrPhone string,
	eventType string,
) (bool, *OTP, error) {

	rdbClient := roo.RDB

	var otp *OTP

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
	otpExists, existErr := rdbClient.Exists(context.Background(), key).Result()

	if existErr != nil {
		return false, nil, fmt.Errorf("error checking existence of OTP: %s", existErr.Error())
	}

	if otpExists != 1 {
		return false, nil, nil
	}

	// user exists and fetching and putting value in User
	val, err := rdbClient.Get(context.Background(), key).Result()
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

func (roo *RedisOTPOperator) Set(
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

	rdbClient := roo.RDB

	// Setting No of minutes as expiration time in Redis
	setErr := rdbClient.Set(
		context.Background(), 
		fmt.Sprintf("%s_%s", emailOrPhone, eventType), 
		encoding.BinaryMarshaler(OTP{
			OtpNumber: otpNumber,
			Event: eventType,
			Type: otpType,
			Email: emailOrPhone,
			Phone: emailOrPhone,
		}),
		(time.Duration(noOfMinutes) * time.Minute),
	).Err()

	if setErr != nil {
		return setErr
	}

	return nil
}