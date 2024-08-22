package service

import (
	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/dao"
)

func GetPRSentByMe(
	input models.PairRequestsSentInput, 
	userId string,
) (bool, models.PairRequestSentByIdOutput) {
	dbUserOperator, dbReceiveErr := dao.GetDBAccessOperator()
	if dbReceiveErr != nil {
		return false, models.PairRequestSentByIdOutput{
			Status: false,
			ErrorMsg: dbReceiveErr.Error(),
		}
	}

	pairReq, fetchPRErr := dbUserOperator.FetchPairRequestsSentByID(
		userId, 
		input.Limit,
		input.Offset,
	)

	if fetchPRErr != nil {
		return false, models.PairRequestSentByIdOutput{
			Status: false,
			ErrorMsg: fetchPRErr.Error(),
		}
	}

	return true, models.PairRequestSentByIdOutput{
		Status: true,
		PairRequests: pairReq,
	}
}

func GetPRSentToMe(
	input models.PairRequestsSentInput, 
	userId string,
) (bool, models.PairRequestSentToIdOutput) {
	dbUserOperator, dbReceiveErr := dao.GetDBAccessOperator()
	if dbReceiveErr != nil {
		return false, models.PairRequestSentToIdOutput{
			Status: false,
			ErrorMsg: dbReceiveErr.Error(),
		}
	}

	pairReq, fetchPRErr := dbUserOperator.FetchPairRequestsSentToID(
		userId, 
		input.Limit,
		input.Offset,
	)

	if fetchPRErr != nil {
		return false, models.PairRequestSentToIdOutput{
			Status: false,
			ErrorMsg: fetchPRErr.Error(),
		}
	}

	return true, models.PairRequestSentToIdOutput{
		Status: true,
		PairRequests: pairReq,
	}
}