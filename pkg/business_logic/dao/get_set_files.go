package dao

import (
	"errors"
	"github.com/drink-events-backend/models"
	"gorm.io/gorm"
)

func (dao *DatabaseAccessOperator) GetFileFromId(
	id string,
) (*models.FileObject, error) {
	var file *models.FileObject
	if fetchErr := dao.DB.Raw("SELECT * from files where id = $1 LIMIT 1;", id).Scan(&file).Error; fetchErr != nil {
		if errors.Is(fetchErr, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, errors.New("error fetching file by id")
	}

	return file, nil
}

func (dao *DatabaseAccessOperator) SetFile(
	f *models.FileObject,
) error {
	dao.Lock()
	defer dao.Unlock()

	if fileAdditionErr := dao.DB.Exec("INSERT INTO files (id, file_name, etag, secure_url, aws_key) values ($1, $2, $3, $4, $5);", f.Id, f.FileName, f.Etag, f.SecureURL, f.AWS_Key).Error; fileAdditionErr != nil {
		return errors.New("error adding file to DB")
	}
	return nil
}

func (dao *DatabaseAccessOperator) DeleteFileByID(id string) error {
	dao.Lock()
	defer dao.Unlock()

	if fileDeletionErr := dao.DB.Exec("DELETE FROM files Where id = ?;", id).Error; fileDeletionErr != nil && !errors.Is(fileDeletionErr, gorm.ErrRecordNotFound) {
		return errors.New("error deleting file")
	}
	return nil
}