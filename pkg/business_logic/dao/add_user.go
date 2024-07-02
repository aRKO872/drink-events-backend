package dao

import "github.com/drink-events-backend/models"

func (dao *DatabaseAccessOperator) AddNormalUser(
	u *models.Users,
) (error) {
	dao.Lock()
	defer dao.Unlock()

	if userAdditionErr := dao.DB.Exec("INSERT INTO users (id, name, user_type, email, phone, bio) values ($1, $2, 'user', $3, $4, $5);", u.Id, u.Name, u.Email, u.Phone, u.Bio).Error; userAdditionErr != nil {
		return userAdditionErr
	}
	return nil
}