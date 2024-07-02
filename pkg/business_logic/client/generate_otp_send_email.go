package client

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/rao"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
)

func GenerateEmailOTPAndSendEmail(
	user *models.Users,
	otpRedisHelper *rao.RedisAccessOperator,
) {
	email := user.Email
	otpNum := rand.Intn(899999) + 100000

	// Save in Redis with 5 mins expiration
	if redisSetErr := otpRedisHelper.SetOTP(
		context.Background(),
		otpNum,
		"email",
		email,
		"verify_email",
		5,
	); redisSetErr != nil {
		log.Println("Here is the error 1, ", redisSetErr)
		return
	}

	// Send Email
	msg := []byte(fmt.Sprintf("To: %s\r\n", email) +
		"Subject: OTP Verification : Drink Events\r\n" +
		"\r\n" +
		fmt.Sprintf("Your OTP is %d\r\n", otpNum))

	to := []string{email}

	log.Println(otpNum, " is sent successfully, or there is this error: ")
	if sendEmailErr := pkg_helpers.SendEmail(to, msg); sendEmailErr != nil {
		fmt.Println(sendEmailErr.Error())
		log.Println("Here is the error 2")
		return
	}
}