package service

import (
	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/dao"
	"github.com/drink-events-backend/pkg/business_logic/rao"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

func DeletePairRequest(
	ctx context.Context,
	pairRequestId string, 
	userId string,
) models.DeletePairRequestResponse {
	// validate if pair request sender or receiver is user.
	// if yes, delete from redis and DB

	dbUserOperator, dbReceiveErr := dao.GetDBAccessOperator()
	if dbReceiveErr != nil {
		return models.DeletePairRequestResponse{
			Status: false,
			ErrorMsg: dbReceiveErr.Error(),
		}
	}

	redisUserOperator, rdbReceiveErr := rao.GetRedisAccessOperator()
	if rdbReceiveErr != nil {
		return models.DeletePairRequestResponse{
			Status: false,
			ErrorMsg: rdbReceiveErr.Error(),
		}
	}

	pairRequest, fetchPRErr := dbUserOperator.GetPairRequestByID(pairRequestId, userId)
	if fetchPRErr != nil {
		return models.DeletePairRequestResponse{
			Status: false,
			ErrorMsg: fetchPRErr.Error(),
		}
	}

	wg := new(errgroup.Group)

	wg.Go(func() error{
		return dbUserOperator.DeletePairRequestByID(pairRequestId)
	})

	wg.Go(func() error{
		return redisUserOperator.DeletePairRequest(ctx, pairRequest.SenderID, pairRequest.ReceiverID)
	})

	if err := wg.Wait(); err != nil {
		return models.DeletePairRequestResponse{
			Status: false,
			ErrorMsg: err.Error(),
		}
	}

	return models.DeletePairRequestResponse{
		Status: true,
		ErrorMsg: "pair request deleted successfully",
	}
}