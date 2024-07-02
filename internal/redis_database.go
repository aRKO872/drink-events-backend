package internal_database

import (
	"fmt"
	"sync"

	pkg_config "github.com/drink-events-backend/pkg/config"
	"github.com/redis/go-redis/v9"
)

var (
	redisOp *RedisOperator
)

func initRDB () (*redis.Client, error) {
	if redisOp != nil {
		return nil, fmt.Errorf("redis database already initialized")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: pkg_config.GetProjectConfig().REDIS_PASSWORD, // no password set
		DB:       0,  // use default DB
	});

	return rdb, nil
}

func GetRDB () (*RedisOperator, error) {
	if redisOp != nil {
		return redisOp, nil;
	}

	// Initialize Redis DB
	rdb, initializeErr := initRDB()

	if initializeErr != nil {
		return nil, initializeErr
	}

	redisOp = &RedisOperator{
		new(sync.RWMutex),
		rdb,
	}
	return redisOp, nil;
}

type RedisOperator struct {
	*sync.RWMutex
	RDB *redis.Client
}

func CreateOrGetRedisOperator() (*RedisOperator, error) {
	rdb, rdbError := GetRDB()
	if rdbError != nil {
		return nil, rdbError
	}

	return rdb, nil
}