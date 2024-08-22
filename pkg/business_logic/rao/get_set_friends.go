package rao

import (
	"context"
	"errors"
	"fmt"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
	"github.com/redis/go-redis/v9"
)

func (rao *RedisAccessOperator) GetFriendRecord(
	ctx context.Context,
	friend1ID, friend2ID string,
) (*models.Friends, error) {
	rdbClient := rao.RDB

	var f models.Friends

	smallId, bigId := pkg_helpers.Smaller(friend1ID, friend2ID)

	// user exists and fetching and putting value in User
	fetchErr := rdbClient.Get(ctx, fmt.Sprintf(literals.FRIEND_RECORD_REDIS_KEY, smallId, bigId)).Scan(&f)

	if fetchErr != nil {
		if errors.Is(fetchErr, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching pair request data: %s", fetchErr.Error())
	}

	return &f, nil
}

func (rao *RedisAccessOperator) AddFriendRecord(
	ctx context.Context,
	f *models.Friends,
) error {
	rdbClient := rao.RDB

	rao.Lock()
	defer rao.Unlock()

	smallId, bigId := pkg_helpers.Smaller(f.Friend1ID, f.Friend2ID)

	setErr := rdbClient.Set(ctx, fmt.Sprintf(literals.FRIEND_RECORD_REDIS_KEY, smallId, bigId), f, 0).Err()

	if setErr != nil {
		return errors.New("error setting friend record in redis")
	}

	return nil
}

func (rao *RedisAccessOperator) DeleteFriendRecord(
	ctx context.Context,
	friend1ID, friend2ID string,
) (error) {
	rdbClient := rao.RDB

	rao.Lock()
	defer rao.Unlock()

	smallId, bigId := pkg_helpers.Smaller(friend1ID, friend2ID)

	if err := rdbClient.Del(ctx, fmt.Sprintf(literals.FRIEND_RECORD_REDIS_KEY, smallId, bigId)); err.Err() != nil {
		if errors.Is(err.Err(), redis.Nil) {
			return errors.New("friend record not found")
		}
		return errors.New("failed to delete friend record from redis")
	}

	return nil
}