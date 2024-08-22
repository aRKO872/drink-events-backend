package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/dao"
	"github.com/drink-events-backend/pkg/business_logic/rao"
)

func GetUser(
	ctx context.Context, 
	id string,
) (*models.Users, error) {
	dbUserOperator, dbReceiveErr := dao.GetDBAccessOperator()
	if dbReceiveErr != nil {
		return nil, dbReceiveErr
	}

	redisUserOperator, rdbReceiveErr := rao.GetRedisAccessOperator()
	if rdbReceiveErr != nil {
		return nil, dbReceiveErr
	}

	fetchStatus, user, fetchUserErr := redisUserOperator.GetUser(ctx, id)

	if fetchUserErr != nil {
		fmt.Println(fetchUserErr)
		return nil, errors.New("error fetching user dtls from redis")
	}

	if !fetchStatus {
		var err error
		// Not available in cache
		user, err = dbUserOperator.GetUserFromId(id)
		if err != nil {
			return nil, errors.New("db error fetching User from DB using ID")
		}

		// Save in RDB before sending
		redisUserOperator.SetUser(ctx, id, user)
	}

	if !user.IsActive {
		return nil, errors.New("fetching data for inactive user")
	}

	return user, nil
}