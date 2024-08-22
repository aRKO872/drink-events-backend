package dao

import (
	"errors"
	"time"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	"gorm.io/gorm"
)

func (dao *DatabaseAccessOperator) AddPairRequest(
	pr *models.PairRequest,
) error {
	dao.Lock()
	defer dao.Unlock()

	if prAdditionErr := dao.DB.Exec("INSERT INTO pair_requests (id, sender_id, receiver_id, is_active, created_at, updated_at) values ($1, $2, $3, $4, $5, $6);", pr.Id, pr.SenderID, pr.ReceiverID, pr.IsActive, pr.CreatedAt, pr.UpdatedAt).Error; prAdditionErr != nil {
		return errors.New("error adding pair request to DB")
	}
	return nil
}

func (dao *DatabaseAccessOperator) DeletePairRequestByID(id string) error {
	dao.Lock()
	defer dao.Unlock()

	updateTime := time.Now().Format(literals.DATE_FORMAT)
	if prDeletionErr := dao.DB.Exec("UPDATE pair_requests SET is_active = false, updated_at = ? Where id = ?;", updateTime, id).Error; prDeletionErr != nil {
		return errors.New("error deleting pair request")
	}

	return nil
}

func (dao *DatabaseAccessOperator) GetPairRequestByID(
	pairRequestId string,
	userId string,
) (models.PairRequest, error) {
	var pairRequest models.PairRequest

	result := dao.DB.Raw(`select * from pair_requests where id = $1 and (
	sender_id = $2 or
	receiver_id = $2
	) and is_active = true limit 1;`, pairRequestId, userId).Scan(&pairRequest)

	if result.Error != nil {
		return models.PairRequest{}, errors.New("error fetching pair request record data")
	}

	if result.RowsAffected == 0 {
		return models.PairRequest{}, errors.New("pair request does not exist")
	}

	return pairRequest, nil
}

func (dao *DatabaseAccessOperator) FetchPairRequestsSentByID(
	id string, 
	limit, offset int,
) ([]models.PairRequestsSentByIdItem, error){
	var pairRequestsSentById []models.PairRequestsSentByIdItem

	dbStr := `select pr.id as id, 
						pr.receiver_id, 
						f.secure_url as profile_pic,
						u.name as receiver_name, 
						pr.created_at
						from pair_requests pr join users u on pr.receiver_id = u.id and u.is_active = true
						left join files f on u.profile_picture = f.id
						where pr.sender_id = $1 
						and pr.is_active = true 
						order by pr.created_at desc Limit $2 offset $3;`
	if fetchErr := dao.DB.Raw(dbStr, id, limit, offset).Scan(&pairRequestsSentById).Error; fetchErr != nil {
		if errors.Is(fetchErr, gorm.ErrRecordNotFound) {
			return []models.PairRequestsSentByIdItem{}, nil
		} else {
			return []models.PairRequestsSentByIdItem{}, errors.New("error fetching records")
		}
	}

	return pairRequestsSentById, nil
}

func (dao *DatabaseAccessOperator) FetchPairRequestsSentToID(
	id string, 
	limit, offset int,
) ([]models.PairRequestsSentToIdItem, error){
	
	var pairRequestsSentToId []models.PairRequestsSentToIdItem
	dbStr := `select pr.id as id, 
						pr.sender_id,
						f.secure_url as profile_pic, 
						u.name as sender_name, 
						pr.created_at
						from pair_requests pr join users u on pr.sender_id = u.id and u.is_active = true
						left join files f on u.profile_picture = f.id
						where pr.receiver_id = $1 
						and pr.is_active = true 
						order by pr.created_at desc Limit $2 offset $3;`
	if fetchErr := dao.DB.Raw(dbStr, id, limit, offset).Scan(&pairRequestsSentToId).Error; fetchErr != nil {
		if errors.Is(fetchErr, gorm.ErrRecordNotFound) {
			return []models.PairRequestsSentToIdItem{}, nil
		} else {
			return []models.PairRequestsSentToIdItem{}, errors.New("error fetching records")
		}
	}

	return pairRequestsSentToId, nil
}