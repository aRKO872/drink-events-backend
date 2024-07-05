package rao

import (
	"context"
	"fmt"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	"github.com/redis/go-redis/v9"
)

func (rao *RedisAccessOperator) SetGEOLocation(
	ctx context.Context,
	locationEvt models.LocationSetEvent, 
	userID string,
) (error) {
	rao.Lock()
	defer rao.Unlock()

	if setOwnLocationErr := rao.UpdateUserLocation(ctx, locationEvt, userID, false); setOwnLocationErr != nil {
		return setOwnLocationErr
	}

	resObj := rao.RDB.GeoAdd(
		ctx, 
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

func (rao *RedisAccessOperator) GetPeopleNearUserGeoLocation(
	ctx context.Context,
	locationEvt models.LocationSetEvent, 
	userId string,
	searchRadius int,
) ([]redis.GeoLocation, error) {
	if err := rao.SetGEOLocation(ctx, locationEvt, userId); err != nil {
		return []redis.GeoLocation{}, fmt.Errorf("error setting geo location for user")
	}

	rao.Lock()
	defer rao.Unlock()

	

	geoLocations, fetchLocationErr := rao.RDB.GeoRadius(
		ctx,
		literals.USER_LOCATION_KEY,
		locationEvt.Longitude,
		locationEvt.Latitude,
		&redis.GeoRadiusQuery{
			Radius: float64(searchRadius),
			Unit: "m",
			WithCoord: true,
			WithDist: true,
			WithGeoHash: true,
		},
	).Result()

	if fetchLocationErr != nil {
		return []redis.GeoLocation{}, fmt.Errorf("error fetching nearby users")
	}

	return geoLocations, nil
}