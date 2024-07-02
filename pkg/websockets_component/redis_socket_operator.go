package pkg_component_websockets

import (
	"context"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	redis_operations "github.com/drink-events-backend/internal"
	"github.com/redis/go-redis/v9"
)

type RedisSocketOperator struct {
	redis_operations.RedisOperator
}

func GetRedisSocketOperator() (*RedisSocketOperator, error) {
	ro, initErr := redis_operations.CreateOrGetRedisOperator()

	if initErr != nil {
		return nil, initErr
	}

	return &RedisSocketOperator{
		RedisOperator: *ro,
	}, nil
}

func (rso *RedisSocketOperator) SetGEOLocation(locationEvt models.LocationSetSocketEvent, userID string) (error) {
	rso.Lock()
	defer rso.Unlock()
	resObj := rso.RDB.GeoAdd(
		context.Background(), 
		literals.USER_LOCATION_KEY,
		&redis.GeoLocation{
			Name: userID,
			Longitude: locationEvt.Longitude,
			Latitude: locationEvt.Latitude,
		},
	)

	if resObj.Err() != nil {
		return resObj.Err()
	}
	return nil
}