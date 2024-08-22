package rao

import (
	"context"
	"errors"
	"fmt"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	"github.com/redis/go-redis/v9"
)

func (rao *RedisAccessOperator) GetFileRecord(
	ctx context.Context,
	id string,
) (*models.FileObject, error) {
	rdbClient := rao.RDB

	var file models.FileObject

	// Checking if exists user exists 
	fileExists, existErr := rdbClient.Exists(ctx, fmt.Sprintf(literals.FILE_INFO_REDIS_KEY, id)).Result()

	if fileExists != 1 || existErr != nil {
		return nil, fmt.Errorf("error checking existence of file: %s", existErr.Error())
	}

	// user exists and fetching and putting value in User
	fetchErr := rdbClient.Get(ctx, fmt.Sprintf(literals.FILE_INFO_REDIS_KEY, id)).Scan(&file)

	if fetchErr != nil {
		return nil, fmt.Errorf("error fetching file data: %s", fetchErr.Error())
	}

	return &file, nil
}

func (rao *RedisAccessOperator) SetFileRecord(
	ctx context.Context,
	id string, 
	file *models.FileObject,
) error {
	rdbClient := rao.RDB

	rao.Lock()
	defer rao.Unlock()

	setErr := rdbClient.Set(ctx, fmt.Sprintf(literals.FILE_INFO_REDIS_KEY, id), file, 0).Err()

	if setErr != nil {
		return setErr
	}

	return nil
}

func (rao *RedisAccessOperator) DeleteFileRecord(
	ctx context.Context,
	id string,
) (error) {
	rao.Lock()
	defer rao.Unlock()

	if err := rao.RDB.Del(ctx, fmt.Sprintf(literals.FILE_INFO_REDIS_KEY, id)); err.Err() != nil && !errors.Is(err.Err(), redis.Nil) {
		return errors.New("failed to delete file from redis")
	}

	return nil
}