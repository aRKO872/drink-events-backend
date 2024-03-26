package pkg_component_user

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	internal_database "github.com/drink-events-backend/internal"
	pkg_component "github.com/drink-events-backend/pkg/components/otp_component"
	pkg_config "github.com/drink-events-backend/pkg/config"
	endpoint_inputs "github.com/drink-events-backend/pkg/endpoint-inputs"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
}

func (user *Users) MarshalBinary() ([]byte, error) {
	return json.Marshal(user)
}

func (u *Users) SignUp() (status bool, output *endpoint_inputs.SignUpLoginOutput) {
	// Check if user exists with Email and phone number provided in DB
	// If exists throw error to login
	// If not exists :
	// - Create user in DB, and get UUID of User
	// - Use that UUID of user to save in Redis
	// - Make a Refresh Token and Access Token enclosing userType, email, phone, user_id, name and return

	databaseInit, getDBErr := internal_database.GetDB(pkg_config.GetProjectConfig().DATABASE_URL)
	if getDBErr != nil {
		return false, &endpoint_inputs.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: getDBErr.Error(),
		}
	}

	redisUserOp, redisInitErr := GetRedisUserOperator()
	if redisInitErr != nil {
		return false, &endpoint_inputs.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: redisInitErr.Error(),
		}
	}

	var user *Users
	if fetchErr := databaseInit.Raw("SELECT * from users where email = $1 OR phone = $2 LIMIT 1;", u.Email, u.Phone).Scan(&user).Error; fetchErr != nil {
		return false, &endpoint_inputs.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: fetchErr.Error(),
		}
	}

	if user != nil {
		return false, &endpoint_inputs.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: "user already exists. please login",
		}
	}

	// Generate a new UUID for the user ID
	u.Id = uuid.NewString()

	// Insert user successfully in DB
	if userAdditionErr := databaseInit.Exec("INSERT INTO users (id, name, user_type, email, phone, bio) values ($1, $2, 'user', $3, $4, $5);", u.Id, u.Name, u.Email, u.Phone, u.Bio).Error; userAdditionErr != nil {
		return false, &endpoint_inputs.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: fmt.Sprintf("error adding user: %s", userAdditionErr.Error()),
		}
	}

	// Insert user successfully in Redis
	if saveInRedisErr := redisUserOp.Set(u.Id, u); saveInRedisErr != nil {
		return false, &endpoint_inputs.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: fmt.Sprintf("error saving user in redis: %s", saveInRedisErr.Error()),
		}
	}

	// Get Access and Refresh token
	// Generate Access and Refresh Tokens
	access_token, accessTokenErr := generateToken(
		pkg_config.GetProjectConfig().JWT_SECRET_KEY,
		time.Duration(pkg_config.GetProjectConfig().ACCESS_TOKEN_EXPIRY),
		u,
	)

	refresh_token, refreshTokenErr := generateToken(
		pkg_config.GetProjectConfig().JWT_SECRET_KEY,
		time.Duration(pkg_config.GetProjectConfig().REFRESH_TOKEN_EXPIRY),
		u,
	)

	if accessTokenErr != nil {
		return false, &endpoint_inputs.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: fmt.Sprintf("error generating access token: %s", accessTokenErr.Error()),
		}
	}

	if refreshTokenErr != nil {
		return false, &endpoint_inputs.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: fmt.Sprintf("error generating refresh token: %s", refreshTokenErr.Error()),
		}
	}

	return true, &endpoint_inputs.SignUpLoginOutput{
		Status:       true,
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}
}

func (u *Users) LogIn(input *endpoint_inputs.LogInInput) (status bool, output *endpoint_inputs.SignUpLoginOutput) {
	// Check if user exists with Email and phone number provided in DB
	// If Not exists throw error asking to Sign Up
	// If Exists -
	// Set DB data in Redis for the same User
	// Create and Send access-refresh token

	databaseInit, getDBErr := internal_database.GetDB(pkg_config.GetProjectConfig().DATABASE_URL)
	if getDBErr != nil {
		return false, &endpoint_inputs.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: getDBErr.Error(),
		}
	}

	redisUserOp, redisInitErr := GetRedisUserOperator()
	if redisInitErr != nil {
		return false, &endpoint_inputs.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: redisInitErr.Error(),
		}
	}

	var user *Users

	if input.LoggedInFrom == "phone" {
		if fetchErr := databaseInit.Raw("SELECT * from users where phone = $1 LIMIT 1;", input.Phone).Scan(&user).Error; fetchErr != nil {
			return false, &endpoint_inputs.SignUpLoginOutput{
				Status:   false,
				ErrorMsg: fetchErr.Error(),
			}
		}
	} else if input.LoggedInFrom == "email" {
		if fetchErr := databaseInit.Raw("SELECT * from users where email = $1 LIMIT 1;", input.Email).Scan(&user).Error; fetchErr != nil {
			return false, &endpoint_inputs.SignUpLoginOutput{
				Status:   false,
				ErrorMsg: fetchErr.Error(),
			}
		}
	} else {
		return false, &endpoint_inputs.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: "please provide proper from type",
		}
	}
	
	if user == nil {
		return false, &endpoint_inputs.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: "user does not exist. please sign up",
		}
	}

	// Setting Main user to Fetched User
	*u = *user

	// Setting User value in Redis as a goroutine
	go redisUserOp.Set(user.Id, user)

	// Get Access and Refresh token
	// Generate Access and Refresh Tokens
	access_token, accessTokenErr := generateToken(
		pkg_config.GetProjectConfig().JWT_SECRET_KEY,
		time.Duration(pkg_config.GetProjectConfig().ACCESS_TOKEN_EXPIRY),
		u,
	)

	refresh_token, refreshTokenErr := generateToken(
		pkg_config.GetProjectConfig().JWT_SECRET_KEY,
		time.Duration(pkg_config.GetProjectConfig().REFRESH_TOKEN_EXPIRY),
		u,
	)

	if accessTokenErr != nil {
		return false, &endpoint_inputs.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: fmt.Sprintf("error generating access token: %s", accessTokenErr.Error()),
		}
	}

	if refreshTokenErr != nil {
		return false, &endpoint_inputs.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: fmt.Sprintf("error generating refresh token: %s", refreshTokenErr.Error()),
		}
	}

	return true, &endpoint_inputs.SignUpLoginOutput{
		Status:       true,
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}
}

func (u *Users) GetUser(id string) (*Users, error) {
	DB, dbReceiveErr := internal_database.GetDB(pkg_config.GetProjectConfig().DATABASE_URL)
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

func (user *Users) verifyOTPForEmail(otp int, event string) (bool, error) {
	email := user.Email

	roo, redisFetchErr := pkg_component.GetRedisOTPOperator()

	if redisFetchErr != nil {
		return false, redisFetchErr
	}

	_, fetchedOTP, otpFetchErr := roo.Get("email", email, event)

	if otpFetchErr != nil {
		return false, otpFetchErr
	}

	if fetchedOTP.OtpNumber != otp {
		return false, fmt.Errorf("wrong otp provided")
	}

	if removeOTPErr := roo.Remove(user.Email, event); removeOTPErr != nil {
		return false, fmt.Errorf("error removing otp : %s", removeOTPErr.Error())
	}

	return true, nil
}

// Function to verify OTP based on events
func (user *Users) VerifyOTP(otp int, event string) (bool, error) {
	switch event {
	case "verify_email":
		return user.verifyOTPForEmail(otp, event)
	case "verify_phone":
		// TODO when otp for phone is set up
		return true, nil
	default:
		return false, fmt.Errorf("please provide proper event")
	}
}

// Function to Resend OTP
func (user *Users) ResendOTPForVerification(input *endpoint_inputs.ResendOTP) (bool, error) {
	if input.Event == "verify_email" && user.Email == "" {
		return false, fmt.Errorf("please provide email for email verification")
	}

	if input.Event == "verify_phone" && user.Phone == "" {
		return false, fmt.Errorf("please provide phone num for phone verification")
	}

	// Delete existing OTP for mail or phone
	// Send new OTP to mail and Phone

	otpRedisHelper, redisErr := pkg_component.GetRedisOTPOperator()

	if redisErr != nil {
		return false, redisErr
	}

	var removeOTPErr error
	if input.Event == "verify_email" {
		removeOTPErr = otpRedisHelper.Remove(user.Email, input.Event)
		go user.generateEmailOTPAndSendEmail(otpRedisHelper)
	}

	if input.Event == "verify_phone" {
		removeOTPErr = otpRedisHelper.Remove(user.Phone, input.Event)

		// TODO - make a goroutine for sending OTP sms

	}

	if removeOTPErr != nil {
		return false, removeOTPErr
	}

	return true, nil
}

// Function to Generate OTP and Send that OTP to user email
func (user *Users) generateEmailOTPAndSendEmail(otpRedisHelper *pkg_component.RedisOTPOperator) {
	email := user.Email
	otpNum := rand.Intn(899999) + 100000

	// Save in Redis with 5 mins expiration
	if redisSetErr := otpRedisHelper.Set(
		otpNum,
		"email",
		email,
		"verify_email",
		5,
	); redisSetErr != nil {
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
}

// Function to Send OTP to provided email successfully
func (user *Users) VerifyEmail() (bool, error) {
	email := user.Email

	if !pkg_helpers.IsValidEmail(email) {
		return false, fmt.Errorf("email invalid")
	}

	// check if an OTP has already been sent for email verification in 5 minute time window? if yes, return false
	otpRedisHelper, redisErr := pkg_component.GetRedisOTPOperator()

	if redisErr != nil {
		return false, redisErr
	}

	otpExists, _, _ := otpRedisHelper.Get(
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
	go user.generateEmailOTPAndSendEmail(otpRedisHelper)

	return true, nil
}
