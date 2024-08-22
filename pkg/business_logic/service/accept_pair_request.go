package service

import (
	"context"
	"errors"
	"time"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/dao"
	"github.com/drink-events-backend/pkg/business_logic/rao"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

func CreateNewFriend(
	ctx context.Context,
	userId string,
	input models.AcceptPairRequestInput,
) (models.Friends, models.PRAcceptOutputWSEvent, error) {
	// Update Pair Request isActive = false in DB
	// Delete Pair Request Object from Redis
	// Create new Friends Object, persist it in DB and Redis

	dbUserOperator, dbReceiveErr := dao.GetDBAccessOperator()
	if dbReceiveErr != nil {
		return models.Friends{}, models.PRAcceptOutputWSEvent{}, dbReceiveErr
	}

	redisUserOperator, rdbReceiveErr := rao.GetRedisAccessOperator()
	if rdbReceiveErr != nil {
		return models.Friends{}, models.PRAcceptOutputWSEvent{}, rdbReceiveErr
	}

	recordId := uuid.NewString()
	creationTime := time.Now().Format(literals.DATE_FORMAT)

	friendRecord := models.Friends{
		Id: recordId,
		Friend1ID: userId,
		Friend2ID: input.SenderId,
		IsActive: true,
		CreatedAt: creationTime,
		UpdatedAt: creationTime,
	}

	prAcceptWSChan := make(chan models.PRAcceptOutputWSEvent, 1)

	wg := new(errgroup.Group)

	wg.Go(func() error{
		deleteOutput := DeletePairRequest(ctx, input.PairRequestId, userId)

		if !deleteOutput.Status {
			return errors.New(deleteOutput.ErrorMsg)
		}
		return nil
	})

	wg.Go(func() error{
		return redisUserOperator.AddFriendRecord(
			ctx, 
			&friendRecord,
		)
	})

	wg.Go(func() error{
		return dbUserOperator.AddFriendRecord(&friendRecord)
	})

	wg.Go(func() error{
		_, user, fetchUserErr := redisUserOperator.GetUser(ctx, userId)

		if fetchUserErr != nil {
			return fetchUserErr
		}

		if user.ProfilePicture == "" {
			prAcceptWSChan <- models.PRAcceptOutputWSEvent{
				AccepterName: user.Name,
				ProfilePicAccepter: "",
			}
	
			return nil
		}

		file, fetchFileErr := redisUserOperator.GetFileRecord(ctx, user.ProfilePicture)

		if fetchFileErr != nil {
			return fetchFileErr
		}

		prAcceptWSChan <- models.PRAcceptOutputWSEvent{
			AccepterName: user.Name,
			ProfilePicAccepter: file.SecureURL,
		}

		return nil
	})

	if err := wg.Wait(); err != nil {
		return models.Friends{}, models.PRAcceptOutputWSEvent{}, err
	}

	close(prAcceptWSChan)

	prAcceptWS := <-prAcceptWSChan
	prAcceptWS.FriendRecordId = friendRecord.Id

	return friendRecord, prAcceptWS, nil
}