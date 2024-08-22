package dao

import (
	"errors"
	"time"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
)

func (dao *DatabaseAccessOperator) GetFriendRecordByFriendsID(
	friend1Id string,
	friend2Id string,
	sendErrorForPresenceOfRecord bool,
) (models.Friends, error) {
	var friend models.Friends

	result := dao.DB.Raw(`select * from friends where ((friend1_id = $1 and friend2_id = $2) or
	(friend1_id = $2 and friend2_id = $1)) and is_active = true limit 1;`, friend1Id, friend2Id).Scan(&friend)

	if result.Error != nil {
		return models.Friends{}, errors.New("error fetching friends record data")
	}

	if result.RowsAffected == 0 {
		return models.Friends{}, errors.New("friends record does not exist")
	} else if sendErrorForPresenceOfRecord {
		return models.Friends{}, errors.New("already friends with the person")
	}

	return friend, nil
}

func (dao *DatabaseAccessOperator) AddFriendRecord(
	f *models.Friends,
) error {
	dao.Lock()
	defer dao.Unlock()

	if friendAdditionErr := dao.DB.Exec("INSERT INTO friends (id, friend1_id, friend2_id, is_active, created_at, updated_at) values ($1, $2, $3, $4, $5, $6);", f.Id, f.Friend1ID, f.Friend2ID, f.IsActive, f.CreatedAt, f.UpdatedAt).Error; friendAdditionErr != nil {
		return errors.New("error adding friend record to DB")
	}
	return nil
}

func (dao *DatabaseAccessOperator) DeleteFriendRecordByID(id string) error {
	dao.Lock()
	defer dao.Unlock()

	updateTime := time.Now().Format(literals.DATE_FORMAT)
	if friendDeletionErr := dao.DB.Exec("UPDATE friends SET is_active = false, updated_at = ? Where id = ?;", updateTime, id).Error; friendDeletionErr != nil {
		return errors.New("error deleting friend record")
	}

	return nil
}