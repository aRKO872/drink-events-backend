package dao

import (
	"errors"
	"time"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	"gorm.io/gorm"
)

func (dao *DatabaseAccessOperator) AddNormalUser(
	u *models.Users,
) (error) {
	dao.Lock()
	defer dao.Unlock()

	u.CreatedAt = time.Now().Format(literals.DATE_FORMAT)
	u.UpdatedAt = u.CreatedAt
	u.SearchRadius = 5000

	if userAdditionErr := dao.DB.Exec("INSERT INTO users (id, name, user_type, email, phone, bio, created_at, updated_at, search_radius) values ($1, $2, 'user', $3, $4, $5, $6, $7, $8);", u.Id, u.Name, u.Email, u.Phone, u.Bio, u.CreatedAt, u.UpdatedAt, u.SearchRadius).Error; userAdditionErr != nil {
		return userAdditionErr
	}
	return nil
}

func (dao *DatabaseAccessOperator) FetchUserFromEmailOrPhone(
	email string,
	phone string,
) (*models.Users, error) {
	var user *models.Users
	if fetchErr := dao.DB.Raw("SELECT * from users where email = $1 OR phone = $2 LIMIT 1;", email, phone).Scan(&user).Error; fetchErr != nil {
		return nil, errors.New("error fetching user by email or phone")
	}

	return user, nil
}

func (dao *DatabaseAccessOperator) FetchUserFromEmail(
	email string,
) (*models.Users, error) {
	var user *models.Users
	if fetchErr := dao.DB.Raw("SELECT * from users where email = $1 LIMIT 1;", email).Scan(&user).Error; fetchErr != nil {
		return nil, errors.New("error fetching user by email")
	}

	return user, nil
}

func (dao *DatabaseAccessOperator) FetchUserFromPhone(
	phone string,
) (*models.Users, error) {
	var user *models.Users
	if fetchErr := dao.DB.Raw("SELECT * from users where phone = $1 LIMIT 1;", phone).Scan(&user).Error; fetchErr != nil {
		return nil, errors.New("error fetching user by phone")
	}

	return user, nil
}

func (dao *DatabaseAccessOperator) GetUserFromId(
	id string,
) (*models.Users, error) {
	var user *models.Users
	if fetchErr := dao.DB.Raw("SELECT * from users where id = $1 LIMIT 1;", id).Scan(&user).Error; fetchErr != nil {
		if errors.Is(fetchErr, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, errors.New("error fetching user by id")
	}

	return user, nil
}