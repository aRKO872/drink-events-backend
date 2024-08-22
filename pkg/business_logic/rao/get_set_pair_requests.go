package rao

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
	"github.com/redis/go-redis/v9"
)

func (rao *RedisAccessOperator) GetPairRequest(
	ctx context.Context,
	senderId, receiverId string,
) (*models.PairRequest, error) {
	rdbClient := rao.RDB

	var pr models.PairRequest

	smallId, bigId := pkg_helpers.Smaller(senderId, receiverId)

	// user exists and fetching and putting value in User
	fetchErr := rdbClient.Get(ctx, fmt.Sprintf(literals.PAIR_REQUEST_REDIS_KEY, smallId, bigId)).Scan(&pr)

	if fetchErr != nil {
		if errors.Is(fetchErr, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching pair request data: %s", fetchErr.Error())
	}

	return &pr, nil
}

func (rao *RedisAccessOperator) AddPairRequest(
	ctx context.Context,
	pr *models.PairRequest,
) error {
	rdbClient := rao.RDB

	rao.Lock()
	defer rao.Unlock()

	smallId, bigId := pkg_helpers.Smaller(pr.SenderID, pr.ReceiverID)

	setErr := rdbClient.Set(ctx, fmt.Sprintf(literals.PAIR_REQUEST_REDIS_KEY, smallId, bigId), pr, 0).Err()

	if setErr != nil {
		return errors.New("failed to add pair request in redis")
	}

	return nil
}

func (rao *RedisAccessOperator) DeletePairRequest(
	ctx context.Context,
	senderID, receiverID string,
) (error) {

	rdbClient := rao.RDB

	rao.Lock()
	defer rao.Unlock()

	smallId, bigId := pkg_helpers.Smaller(senderID, receiverID)

	log.Println(smallId, bigId)

	err := rdbClient.Del(ctx, fmt.Sprintf(literals.PAIR_REQUEST_REDIS_KEY, smallId, bigId)); 
	if err.Err() != nil {
		if errors.Is(err.Err(), redis.Nil) {
			return errors.New("pair request not found")
		}
		return errors.New("failed to delete pair request from redis")
	}

	log.Println(err)

	return nil
}