package internal_database

import (
	"fmt"

	pkg_config "github.com/drink-events-backend/pkg/config"
	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
)

func initRDB () error {
	if rdb != nil {
		return fmt.Errorf("redis database already initialized")
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: pkg_config.GetProjectConfig().REDIS_PASSWORD, // no password set
		DB:       0,  // use default DB
	});

	return nil
}

func GetRDB () (*redis.Client, error) {
	if rdb != nil {
		return rdb, nil;
	}

	// Initialize Redis DB
	initializeErr := initRDB()

	if initializeErr != nil {
		return nil, initializeErr
	}

	return rdb, nil;
}