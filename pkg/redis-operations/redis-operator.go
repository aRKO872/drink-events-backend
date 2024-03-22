package redis_operations

import (
	"github.com/drink-events-backend/internal"
	"github.com/redis/go-redis/v9"
)

type RedisOperator struct {
	RDB *redis.Client
}

func CreateOrGetRedisOperator() (*RedisOperator, error) {
	rdb, rdbError := internal_database.GetRDB()
	if rdbError != nil {
		return nil, rdbError
	}

	return &RedisOperator{
		RDB: rdb,
	}, nil
}