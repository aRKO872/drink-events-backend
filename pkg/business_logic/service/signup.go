package service

import (
	"context"
	"fmt"
	"time"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/rao"
	"github.com/drink-events-backend/pkg/business_logic/dao"
	pkg_config "github.com/drink-events-backend/pkg/config"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
	"github.com/google/uuid"
)

func SignUp(
	ctx context.Context,
	u *models.Users,
) (
	status bool, 
	output *models.SignUpLoginOutput,
) {
	// Check if user exists with Email and phone number provided in DB
	// If exists throw error to login
	// If not exists :
	// - Create user in DB, and get UUID of User
	// - Use that UUID of user to save in Redis
	// - Make a Refresh Token and Access Token enclosing userType, email, phone, user_id, name and return

	dao, getDBErr := dao.GetDBAccessOperator()
	if getDBErr != nil {
		return false, &models.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: getDBErr.Error(),
		}
	}

	redisUserOp, redisInitErr := rao.GetRedisAccessOperator()
	if redisInitErr != nil {
		return false, &models.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: redisInitErr.Error(),
		}
	}

	var user *models.Users
	var fetchErr error

	user, fetchErr = dao.FetchUserFromEmailOrPhone(u.Email, u.Phone); 
	if fetchErr != nil {
		return false, &models.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: fetchErr.Error(),
		}
	}

	if user != nil {
		return false, &models.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: "user already exists. please login",
		}
	}

	// Generate a new UUID for the user ID
	u.Id = uuid.NewString()

	// Insert user successfully in DB
	if userAdditionErr := dao.AddNormalUser(u); userAdditionErr != nil {
		return false, &models.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: fmt.Sprintf("error adding user: %s", userAdditionErr.Error()),
		}
	}

	// Insert user successfully in Redis
	if saveInRedisErr := redisUserOp.SetUser(ctx, u.Id, u); saveInRedisErr != nil {
		return false, &models.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: fmt.Sprintf("error saving user in redis: %s", saveInRedisErr.Error()),
		}
	}

	// Get Access and Refresh token
	// Generate Access and Refresh Tokens
	access_token, accessTokenErr := pkg_helpers.GenerateToken(
		pkg_config.GetProjectConfig().JWT_SECRET_KEY,
		time.Duration(pkg_config.GetProjectConfig().ACCESS_TOKEN_EXPIRY),
		u,
	)

	refresh_token, refreshTokenErr := pkg_helpers.GenerateToken(
		pkg_config.GetProjectConfig().JWT_SECRET_KEY,
		time.Duration(pkg_config.GetProjectConfig().REFRESH_TOKEN_EXPIRY),
		u,
	)

	if accessTokenErr != nil {
		return false, &models.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: fmt.Sprintf("error generating access token: %s", accessTokenErr.Error()),
		}
	}

	if refreshTokenErr != nil {
		return false, &models.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: fmt.Sprintf("error generating refresh token: %s", refreshTokenErr.Error()),
		}
	}

	return true, &models.SignUpLoginOutput{
		Status:       true,
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}
}