package pkg_component

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

func (ruo *RedisUserOperator) Get(key string) (bool, *Users, error) {
	rdbClient := ruo.RDB

	var user *Users

	// Checking if exists user exists 
	userExists, existErr := rdbClient.Exists(context.Background(), key).Result()

	if userExists != 1 || existErr != nil {
		return false, nil, fmt.Errorf("error checking existence of user: %s", existErr.Error())
	}

	// user exists and fetching and putting value in User
	fetchErr := rdbClient.Get(context.Background(), key).Scan(&user)

	if fetchErr != nil {
		return false, nil, fmt.Errorf("error fetching user data: %s", fetchErr.Error())
	}

	return true, user, nil
}

func (ruo *RedisUserOperator) Set(key string, user *Users) error {
	rdbClient := ruo.RDB

	setErr := rdbClient.Set(context.Background(), fmt.Sprintf("%s_user_info", key), user, 0).Err()

	if setErr != nil {
		return setErr
	}

	return nil
}

// CREATE TABLE otps (
//   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
//   number INT NOT NULL CHECK (number >= 100000 AND number <= 999999), 
//   event otp_event_types NOT NULL,
//   expiry TIMESTAMP NOT NULL,
//   created_at TIMESTAMP DEFAULT NOW() NOT NULL,
//   updated_at TIMESTAMP DEFAULT NOW() NOT NULL
// );