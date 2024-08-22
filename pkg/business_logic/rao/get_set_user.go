package rao

import (
	"context"
	"errors"
	"fmt"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	"github.com/redis/go-redis/v9"
)

func (rao *RedisAccessOperator) GetUser(
	ctx context.Context, 
	id string,
) (bool, *models.Users, error) {
	rdbClient := rao.RDB

	var user models.Users

	// Checking if exists user exists 
	userExists, existErr := rdbClient.Exists(ctx, fmt.Sprintf(literals.USER_INFO_REDIS_KEY, id)).Result()

	if userExists != 1 || existErr != nil {
		return false, nil, fmt.Errorf("error checking existence of user: %s", existErr.Error())
	}

	// user exists and fetching and putting value in User
	fetchErr := rdbClient.Get(ctx, fmt.Sprintf(literals.USER_INFO_REDIS_KEY, id)).Scan(&user)

	if fetchErr != nil {
		return false, nil, fmt.Errorf("error fetching user data: %s", fetchErr.Error())
	}

	return true, &user, nil
}

func (ruo *RedisAccessOperator) DeleteUserRecord(
	ctx context.Context,
	user *models.Users,
) (error) {
	if err := ruo.RDB.Del(ctx, fmt.Sprintf(literals.USER_INFO_REDIS_KEY, user.Id)); err.Err() != nil && !errors.Is(err.Err(), redis.Nil) {
		return errors.New("failed to delete user from redis")
	}

	return nil
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