package websockets

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/drink-events-backend/cmd/queue"
	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/dao"
	"github.com/drink-events-backend/pkg/business_logic/rao"
	"github.com/drink-events-backend/pkg/business_logic/service"
	"github.com/gocraft/work"
)

func SendPairRequest(
	evt Event,
	c *Client,
) error {
	var input models.PairRequestWSEvent
	if err := json.Unmarshal(evt.Payload, &input); err != nil {
		return errors.New("error binding input")
	}

	dbUserOperator, dbReceiveErr := dao.GetDBAccessOperator()
	if dbReceiveErr != nil {
		return dbReceiveErr
	}

	redisUserOp, redisInitErr := rao.GetRedisAccessOperator()
	if redisInitErr != nil {
		return redisInitErr
	}

	// Validation for checking not to send friend request if friend link is already there.
	if _, err := dbUserOperator.GetFriendRecordByFriendsID(
		input.To, 
		c.UserId,
		true,
	); err != nil {
		return err
	}

	pr, fetchPRErr := redisUserOp.GetPairRequest(context.Background(), c.UserId, input.To)

	if fetchPRErr != nil {
		return fetchPRErr
	}

	if pr != nil {
		if pr.SenderID == c.UserId {
			return errors.New("you already sent a pair request to them")
		} else {
			return errors.New("they have sent you a pair request previously")
		}
	}

	if toClient, ok := c.Manager.ClientMap[input.To]; ok {
		sendingUser, fetchUserErr := service.GetUser(context.Background(), c.UserId)
		if fetchUserErr != nil {
			return fetchUserErr
		}

		msgPayload := models.PairRequestOutputWSEvent{
			FromName: sendingUser.Name,
			ProfilePic: sendingUser.ProfilePicture,
		}

		payload, err := json.Marshal(msgPayload)
		if err != nil {
			return err
		}
		
		toClient.Feed <- Event {
			Type: literals.SEND_PAIR_REQUEST_SUCCESS,
			Payload: payload,
		}
	}

	_, err := queue.Enqueuer.Enqueue(literals.ADD_PAIR_REQUEST_JOB, work.Q{
		"to": input.To,
		"from": c.UserId,
	})

	if err != nil {
		log.Println(err)
		return errors.New("error persisting pair request info")
	}

	return nil
}