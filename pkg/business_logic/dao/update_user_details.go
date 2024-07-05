package dao

import (
	"errors"
	"time"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
)

func (dao *DatabaseAccessOperator) UpdateUserLocation(
	userID string,
	loc models.LocationSetEvent,
) (error) {
	dao.Lock()
	defer dao.Unlock()

	updateTime := time.Now().Format(literals.DATE_FORMAT)
	if updateErr := dao.DB.Exec("UPDATE users SET latitude = ?, longitude = ?, updated_at = ? WHERE id = ?;", loc.Latitude, loc.Longitude, updateTime, userID).Error; updateErr != nil {
		return errors.New("error updating user location in DB")
	}

	return nil
}