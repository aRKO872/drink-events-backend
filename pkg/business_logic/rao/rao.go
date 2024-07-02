package rao

import (
	redis_operations "github.com/drink-events-backend/internal"
)

type RedisAccessOperator struct {
	redis_operations.RedisOperator
}

func GetRedisAccessOperator() (*RedisAccessOperator, error) {
	ro, initErr := redis_operations.CreateOrGetRedisOperator()

	if initErr != nil {
		return nil, initErr
	}

	return &RedisAccessOperator{
		RedisOperator: *ro,
	}, nil
}