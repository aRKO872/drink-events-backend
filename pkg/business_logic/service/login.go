package service

import (
	"context"
	"fmt"
	"time"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/dao"
	"github.com/drink-events-backend/pkg/business_logic/rao"
	pkg_config "github.com/drink-events-backend/pkg/config"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
)

func LogIn(
	ctx context.Context,
	input *models.LogInInput,
	u *models.Users,
) (status bool, output *models.SignUpLoginOutput) {
	// Check if user exists with Email and phone number provided in DB
	// If Not exists throw error asking to Sign Up
	// If Exists -
	// Set DB data in Redis for the same User
	// Create and Send access-refresh token
	daoUserOp, daoInitErr := dao.GetDBAccessOperator()
	if daoInitErr != nil {
		return false, &models.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: daoInitErr.Error(),
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

	if input.LoggedInFrom == "phone" {
		user, fetchErr = daoUserOp.FetchUserFromPhone(input.Phone)
		if fetchErr != nil {
			return false, &models.SignUpLoginOutput{
				Status:   false,
				ErrorMsg: fetchErr.Error(),
			}
		}
	} else if input.LoggedInFrom == "email" {
		user, fetchErr = daoUserOp.FetchUserFromEmail(input.Email)
		if fetchErr != nil {
			return false, &models.SignUpLoginOutput{
				Status:   false,
				ErrorMsg: fetchErr.Error(),
			}
		}
	} else {
		return false, &models.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: "please provide proper from type",
		}
	}

	if user == nil {
		return false, &models.SignUpLoginOutput{
			Status:   false,
			ErrorMsg: "user does not exist. please sign up",
		}
	}

	// Setting User value in Redis as a goroutine
	go redisUserOp.SetUser(ctx, user.Id, user)

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