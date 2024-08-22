package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/drink-events-backend/models"
	pkg_config "github.com/drink-events-backend/pkg/config"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
)

func FetchUserFromToken(
	r *http.Request,
	sendOnlyClaims bool,
) (*models.Users, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil, fmt.Errorf("missing authorization header")
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	userToken := parts[1]
	userClaims, userFetchedFromTokenErr := pkg_helpers.ParseToken(userToken, pkg_config.GetProjectConfig().JWT_SECRET_KEY)

	if userFetchedFromTokenErr != nil {
		return nil, userFetchedFromTokenErr
	}

	if sendOnlyClaims {
		return &models.Users{
			UserType: userClaims.UserType,
			Id: userClaims.UserId,
		}, nil
	}

	fetchedUser, fetchUserFromUserIDError := GetUser(
		r.Context(), 
		userClaims.UserId,
	)

	if fetchUserFromUserIDError != nil {
		return nil, fetchUserFromUserIDError
	}

	return fetchedUser, nil
}