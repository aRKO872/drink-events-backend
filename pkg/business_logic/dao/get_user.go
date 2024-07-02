package dao

import (
	"errors"
	"fmt"

	"github.com/drink-events-backend/models"
	"gorm.io/gorm"
)

func (dao *DatabaseAccessOperator) FetchUserFromEmailOrPhone(
	email string,
	phone string,
) (*models.Users, error) {
	dao.Lock()
	defer dao.Unlock()

	var user *models.Users
	if fetchErr := dao.DB.Raw("SELECT * from users where email = $1 OR phone = $2 LIMIT 1;", email, phone).Scan(&user).Error; fetchErr != nil {
		return nil, errors.New("error fetching user by email or phone")
	}

	return user, nil
}

func (dao *DatabaseAccessOperator) FetchUserFromEmail(
	email string,
) (*models.Users, error) {
	fmt.Println("past error")
	dao.Lock()
	fmt.Println("post error")
	defer dao.Unlock()

	var user *models.Users
	if fetchErr := dao.DB.Raw("SELECT * from users where email = $1 LIMIT 1;", email).Scan(&user).Error; fetchErr != nil {
		return nil, errors.New("error fetching user by email")
	}

	return user, nil
}

func (dao *DatabaseAccessOperator) FetchUserFromPhone(
	phone string,
) (*models.Users, error) {
	dao.Lock()
	defer dao.Unlock()

	var user *models.Users
	if fetchErr := dao.DB.Raw("SELECT * from users where phone = $1 LIMIT 1;", phone).Scan(&user).Error; fetchErr != nil {
		return nil, errors.New("error fetching user by phone")
	}

	return user, nil
}

func (dao *DatabaseAccessOperator) GetUserFromId(
	id string,
) (*models.Users, error) {
	dao.Lock()
	defer dao.Unlock()

	var user *models.Users
	if fetchErr := dao.DB.Raw("SELECT * from users where id = $1 LIMIT 1;", id).Scan(&user).Error; fetchErr != nil {
		if errors.Is(fetchErr, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, errors.New("error fetching user by id")
	}

	return user, nil
}