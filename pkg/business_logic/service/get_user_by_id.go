package service

import (
	"context"

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

	fetchStatus, user, _ := redisUserOperator.GetUser(ctx, id)

	if fetchStatus {
		// Available in cache
		return user, nil

	} else {
		// Not available in cache
		dbUserOperator.GetUserFromId(id)

		// Save in RDB before sending
		redisUserOperator.SetUser(ctx, id, user)

		return user, nil
	}
}