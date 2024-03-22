package pkg_component

import (
	"errors"
	"fmt"
	"math/rand"

	internal_database "github.com/drink-events-backend/internal"
	pkg_component "github.com/drink-events-backend/pkg/components/otp_component"
	pkg_config "github.com/drink-events-backend/pkg/config"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
	"gorm.io/gorm"
)

type Users struct {
	Id              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	UserType        string    `json:"user_type" db:"user_type"`
	Email           string    `json:"email" db:"email"`
	Phone           string    `json:"phone" db:"phone"`
	ProfilePicture  string    `json:"profile_picture" db:"profile_picture"`
	Bio             string    `json:"bio" db:"bio"`
	CreatedAt       string 		`json:"created_at" db:"created_at"`
	UpdatedAt       string 		`json:"updated_at" db:"updated_at"`
	Latitude        float64   `json:"latitude" db:"latitude"`
	Longitude       float64   `json:"longitude" db:"longitude"`
	EmailLastChanged string   `json:"email_last_changed" db:"email_last_changed"`
	PhoneLastChanged string   `json:"phone_last_changed" db:"phone_last_changed"`
	SearchRadius    int       `json:"search_radius" db:"search_radius"`
}

func (u *Users) GetUser(id string) (*Users, error) {
	DB, dbReceiveErr := internal_database.GetDB(pkg_config.GetProjectConfig().DATABASE_URL);
	if dbReceiveErr != nil {
		return nil, dbReceiveErr
	}

	redisUserOperator, rdbReceiveErr := GetRedisUserOperator()
	if rdbReceiveErr != nil {
		return nil, dbReceiveErr
	}

	fetchStatus, user, _ := redisUserOperator.Get(id)

	
	if fetchStatus {
		// Available in cache
		return user, nil

	} else {
		// Not available in cache
		query := DB.Where("id = ?", id).First(user)

		if query.Error != nil && errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		} else if err := query.Error; err != nil {
				return nil, err
		}

		// Save in RDB before sending
		redisUserOperator.Set(id, user)

		return user, nil
	}
}

// Function to Send OTP to provided email successfully
func (user *Users) VerifyEmail() (bool, error) {
	email := user.Email

	if !pkg_helpers.IsValidEmail(email) {
		return false, fmt.Errorf("email invalid")
	}

	// check if an OTP has already been sent for email verification in 5 minute time window? if yes, return false
	otpRedisHelper, redisErr := pkg_component.GetRedisOTPOperator()

	if redisErr != nil {return false, redisErr}
	
	otpExists, _, _ := otpRedisHelper.Get(
		"email",
		email,
		"verify_email",
	)

	if otpExists {return false, fmt.Errorf("otp already exists in 5 minute time window")}

	// Do below in one Goroutine
	// Generate a new OTP and save it as key `"email"_"event-name"`
	// Send OTP to user's email.
	go func () {
		otpNum := rand.Intn(899999) + 100000
		fmt.Println(otpNum)

		// Save in Redis with 5 mins expiration
		if redisSetErr := otpRedisHelper.Set(
			otpNum,
			"email",
			email,
			"verify_email",
			5,
		); redisSetErr != nil {
			fmt.Println("Here is the error : ", redisSetErr.Error())
			return
		}

		// Send Email
		msg := []byte(fmt.Sprintf("To: %s\r\n", email) +
		"Subject: OTP Verification : Drink Events\r\n" +
		"\r\n" +
		fmt.Sprintf("Your OTP is %d\r\n", otpNum))

		to := []string{email}

		if sendEmailErr := pkg_helpers.SendEmail(to, msg); sendEmailErr != nil {
			fmt.Println(sendEmailErr.Error())
			return
		}
	}();
	
	return true, nil;
}