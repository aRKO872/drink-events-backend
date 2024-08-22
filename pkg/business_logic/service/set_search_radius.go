package service

import (
	"context"
	"time"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/dao"
	"github.com/drink-events-backend/pkg/business_logic/rao"
	"golang.org/x/sync/errgroup"
)

func SetSearchRadius(
	ctx context.Context,
	userId string,
	searchRadius int,
) models.SetUserSearchRadiusResponse {
	userChan := make(chan models.Users, 1)

	dao, getDBErr := dao.GetDBAccessOperator()
	if getDBErr != nil {
		return models.SetUserSearchRadiusResponse{
			Status: false,
			ErrorMsg: getDBErr.Error(),
		}
	}

	redisUserOp, redisInitErr := rao.GetRedisAccessOperator()
	if redisInitErr != nil {
		return models.SetUserSearchRadiusResponse{
			Status: false,
			ErrorMsg: redisInitErr.Error(),
		}
	}

	updateTime := time.Now().Format(literals.DATE_FORMAT)

	wg := new(errgroup.Group)

	wg.Go(func() error {
		_, user, fetchUserErr := redisUserOp.GetUser(ctx, userId)
		if fetchUserErr != nil {
			return fetchUserErr
		}

		user.SearchRadius = searchRadius
		user.UpdatedAt = updateTime
		userChan <- *user

		return redisUserOp.SetUser(ctx, userId, user)
	})

	wg.Go(func() error {
		return dao.SetUserSearchRedius(userId, searchRadius, updateTime)
	})

	if err := wg.Wait(); err != nil {
		return models.SetUserSearchRadiusResponse{
			Status: false,
			ErrorMsg: err.Error(),
		}
	}

	close(userChan)
	user := <-userChan

	return models.SetUserSearchRadiusResponse{
		Status: true,
		ErrorMsg: "successfully change search radius",
		UserData: models.Users{
			Name: user.Name,
			Email: user.Email,
			Phone: user.Phone,
			Latitude: user.Latitude,
			Longitude: user.Longitude,
			SearchRadius: user.SearchRadius,
			UpdatedAt: user.UpdatedAt,
		},
	}
}