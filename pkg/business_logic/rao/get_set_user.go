package rao

import (
	"context"
	"fmt"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
)

func (rao *RedisAccessOperator) GetUser(
	ctx context.Context, 
	id string,
) (bool, *models.Users, error) {
	rdbClient := rao.RDB

	var user *models.Users

	// Checking if exists user exists 
	userExists, existErr := rdbClient.Exists(context.Background(), fmt.Sprintf(literals.USER_INFO_REDIS_KEY, id)).Result()

	if userExists != 1 || existErr != nil {
		return false, nil, fmt.Errorf("error checking existence of user: %s", existErr.Error())
	}

	// user exists and fetching and putting value in User
	rao.Lock()
	defer rao.Unlock()
	fetchErr := rdbClient.Get(ctx, fmt.Sprintf(literals.USER_INFO_REDIS_KEY, id)).Scan(&user)

	if fetchErr != nil {
		return false, nil, fmt.Errorf("error fetching user data: %s", fetchErr.Error())
	}

	return true, user, nil
}

func (ruo *RedisAccessOperator) SetUser(
	ctx context.Context,
	id string, 
	user *models.Users,
) error {
	rdbClient := ruo.RDB

	ruo.Lock()
	defer ruo.Unlock()
	setErr := rdbClient.Set(ctx, fmt.Sprintf(literals.USER_INFO_REDIS_KEY, id), user, 0).Err()

	if setErr != nil {
		return setErr
	}

	return nil
}