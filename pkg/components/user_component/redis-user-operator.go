package pkg_component_user

import (
	"context"
	"fmt"

	redis_operations "github.com/drink-events-backend/pkg/redis-operations"
)

type RedisUserOperator struct {
	redis_operations.RedisOperator
}

func GetRedisUserOperator() (*RedisUserOperator, error) {
	ro, initErr := redis_operations.CreateOrGetRedisOperator()

	if initErr != nil {
		return nil, initErr
	}

	return &RedisUserOperator{
		RedisOperator: *ro,
	}, nil
}

func (ruo *RedisUserOperator) Get(id string) (bool, *Users, error) {
	rdbClient := ruo.RDB

	var user *Users

	// Checking if exists user exists 
	userExists, existErr := rdbClient.Exists(context.Background(), fmt.Sprintf("%s_user_info", id)).Result()

	if userExists != 1 || existErr != nil {
		return false, nil, fmt.Errorf("error checking existence of user: %s", existErr.Error())
	}

	// user exists and fetching and putting value in User
	fetchErr := rdbClient.Get(context.Background(), fmt.Sprintf("%s_user_info", id)).Scan(&user)

	if fetchErr != nil {
		return false, nil, fmt.Errorf("error fetching user data: %s", fetchErr.Error())
	}

	return true, user, nil
}

func (ruo *RedisUserOperator) Set(id string, user *Users) error {
	rdbClient := ruo.RDB

	setErr := rdbClient.Set(context.Background(), fmt.Sprintf("%s_user_info", id), user, 0).Err()

	if setErr != nil {
		return setErr
	}

	return nil
}