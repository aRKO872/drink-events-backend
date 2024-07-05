package rao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
)

func (rao *RedisAccessOperator) UpdateUserLocation(
	ctx context.Context,
	locationObj models.LocationSetEvent,
	userId string,
	standalone bool,
) error {
	if standalone {
		rao.Lock()
		defer rao.Unlock()
	}

	var user models.Users
	fetchErr := rao.RDB.Get(ctx, fmt.Sprintf(literals.USER_INFO_REDIS_KEY, userId)).Scan(&user)

	if fetchErr != nil {
		return errors.New("error fetching user data")
	}

	user.Latitude = locationObj.Latitude
	user.Longitude = locationObj.Longitude
	user.UpdatedAt = time.Now().Format(literals.DATE_FORMAT)

	if setErr := rao.RDB.Set(ctx, fmt.Sprintf(literals.USER_INFO_REDIS_KEY, userId), user, 0).Err(); setErr != nil {
		return setErr
	}

	return nil
}