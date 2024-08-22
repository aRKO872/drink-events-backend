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

	// Save in Redis and fetch nearby userIDs for tokenUser.uID and user.SearchRadius
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

	var userData *models.Users
	g1.Go(func() error {
		var userFetchErr error
		userData, userFetchErr = GetUser(ctx, tokenUser.UserId)
		if userFetchErr != nil {
			return userFetchErr
		}
		geoDataArr, err := redisUserOp.GetPeopleNearUserGeoLocation(ctx, *loc, tokenUser.UserId, userData.SearchRadius)

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

	userDetailsChan := make(chan models.Users, 10)

	for ind := 0; ind < len(nearbyUsersRawArr); ind += 10 {
		end := ind + 10
		if end > len(nearbyUsersRawArr) {
			end = len(nearbyUsersRawArr)
		}

		nearbyUserBatch := nearbyUsersRawArr[ind:end]
		for _, geoUserLoc := range nearbyUserBatch {
			g2.Go(func() error {
				fetchedUser, fetchUserErr := GetUser(ctx, geoUserLoc.Name)
				if fetchUserErr != nil {
					return fetchUserErr
				}
	
				userDetailsChan <- *fetchedUser
				return nil
			})
		}
	}

	if g2Err := g2.Wait(); g2Err != nil {
		return false, &models.GetNearbyUsersOutput{
			Status:   false,
			ErrorMsg: g2Err.Error(),
		}
	}
	close(userDetailsChan)

	var outputNearbyUserArr []models.UsersSelfTag
	for nearbyUser := range userDetailsChan{
		outputNearbyUserArr = append(outputNearbyUserArr, models.UsersSelfTag{
			Users: nearbyUser,
			Self: nearbyUser.Id == tokenUser.UserId,
		})
	}

	return true, &models.GetNearbyUsersOutput{
		Status: true,
		NearbyUsers: outputNearbyUserArr,
	}
}