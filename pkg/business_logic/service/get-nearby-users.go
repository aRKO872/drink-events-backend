package service

import (
	"net/http"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/dao"
	"github.com/drink-events-backend/pkg/business_logic/rao"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

func GetPeopleNearby(
	r *http.Request,
	loc *models.LocationSetEvent,
	tokenUser *models.TokenizedUserDetails,
) (bool, *models.GetNearbyUsersOutput) {
	ctx := r.Context()

	dao, getDBErr := dao.GetDBAccessOperator()
	if getDBErr != nil {
		return false, &models.GetNearbyUsersOutput{
			Status:   false,
			ErrorMsg: getDBErr.Error(),
		}
	}

	redisUserOp, redisInitErr := rao.GetRedisAccessOperator()
	if redisInitErr != nil {
		return false, &models.GetNearbyUsersOutput{
			Status:   false,
			ErrorMsg: redisInitErr.Error(),
		}
	}

	// Save in Redis and fetch nearby userIDs for tokenUser.uID and tokenUser.SearchRadius
	// Set Location for user in Database
	// Create another waitgroup to concurrently fetch all User Details for nearby User IDs.

	g1 := new(errgroup.Group)

	g1.Go(func() error {
		if err := dao.UpdateUserLocation(tokenUser.UserId, *loc); err != nil {
			return err
		}

		return nil
	})

	geoLocArrChan := make(chan []redis.GeoLocation, 1)
	g1.Go(func() error {
		geoDataArr, err := redisUserOp.GetPeopleNearUserGeoLocation(ctx, *loc, tokenUser.UserId, tokenUser.SearchRadius)

		if err != nil {
			return err
		}

		geoLocArrChan <- geoDataArr

		return nil
	})

	if g1Err := g1.Wait(); g1Err != nil {
		return false, &models.GetNearbyUsersOutput{
			Status:   false,
			ErrorMsg: g1Err.Error(),
		}
	}
	close(geoLocArrChan)
	nearbyUsersRawArr := <-geoLocArrChan

	g2 := new(errgroup.Group)

	userDetailsChan := make(chan models.Users, len(nearbyUsersRawArr))

	for _, geoUserLoc := range nearbyUsersRawArr {
		g2.Go(func() error {
			fetchedUser, fetchUserErr := GetUser(ctx, geoUserLoc.Name)
			if fetchUserErr != nil {
				return fetchUserErr
			}

			userDetailsChan <- *fetchedUser
			return nil
		})
	}

	if g2Err := g2.Wait(); g2Err != nil {
		return false, &models.GetNearbyUsersOutput{
			Status:   false,
			ErrorMsg: g2Err.Error(),
		}
	}
	close(userDetailsChan)

	var outputNearbyUserArr []models.Users
	for nearbyUser := range userDetailsChan{
		outputNearbyUserArr = append(outputNearbyUserArr, nearbyUser)
	}

	return true, &models.GetNearbyUsersOutput{
		Status: true,
		NearbyUsers: outputNearbyUserArr,
	}
}