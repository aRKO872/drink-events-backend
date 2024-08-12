package service

import (
	"context"
	"mime/multipart"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/client"
	"github.com/drink-events-backend/pkg/business_logic/dao"
	"github.com/drink-events-backend/pkg/business_logic/rao"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

func saveImageAWSRedisDB(
	ctx context.Context,
	user *models.Users,
	imageFile multipart.File, 
	imageDet *multipart.FileHeader,
	dao *dao.DatabaseAccessOperator, 
	redisUserOp *rao.RedisAccessOperator,
) (models.FileObject, error) {
	fileObject, err := client.SaveObjectToAWS(user.ProfilePicture, imageFile, imageDet)
	if err != nil {
		return models.FileObject{}, err
	}

	// Save file in Redis and DB
	// and update user file id

	wg := new(errgroup.Group)

	wg.Go(func() error {
		if err := redisUserOp.SetFileRecord(ctx, user.ProfilePicture, &fileObject); err != nil {
			return err
		}

		return nil
	})
	
	wg.Go(func() error {
		if err := dao.SetFile(&fileObject); err != nil {
			return err
		}

		if err := dao.SetUserProfileImage(user.Id, user.ProfilePicture); err != nil {
			return err
		}

		return nil
	})

	wg.Go(func() error {
		if err := redisUserOp.SetUser(ctx, user.Id, user); err != nil {
			return err
		}

		return nil
	})

	if err := wg.Wait(); err != nil {
		return models.FileObject{}, err
	}

	return fileObject, nil
}

func removeProfilePictureForUser(
	ctx context.Context, 
	user *models.Users, 
	dao *dao.DatabaseAccessOperator, 
	redisUserOp *rao.RedisAccessOperator,
) error {
	wg := new(errgroup.Group)

	wg.Go(func() error {
		return client.DeleteObjectFromAWS(user.ProfilePicture)
	})

	wg.Go(func() error {
		if err := dao.RemoveUserProfileImage(user.Id); err != nil {
			return err
		}
		return dao.DeleteFileByID(user.ProfilePicture)
	})

	wg.Go(func() error {
		return redisUserOp.DeleteFileRecord(ctx, user.ProfilePicture)
	})

	if err := wg.Wait(); err != nil {
		return err
	}
	return nil
}

func SaveProfilePicture(
	ctx context.Context,
	userDet *models.TokenizedUserDetails,
	imageFile multipart.File,
	imageDet *multipart.FileHeader,
) *models.SavedImageResponse {
	// Fetch User. Check if profile_pic exists. If yes,
	// delete the picture from DB, from redis and from AWS

	// Save to AWS. Get File Object f
	// Edit User Profile Pic img URL in redis and DB 
	// Send back response

	dao, getDBErr := dao.GetDBAccessOperator()
	if getDBErr != nil {
		return &models.SavedImageResponse{
			Status: false,
			ErrorMsg: getDBErr.Error(),
		}
	}

	redisUserOp, redisInitErr := rao.GetRedisAccessOperator()
	if redisInitErr != nil {
		return &models.SavedImageResponse{
			Status: false,
			ErrorMsg: redisInitErr.Error(),
		}
	}

	fetchedUser, fetchUserErr := GetUser(ctx, userDet.UserId)
	if fetchUserErr != nil {
		return &models.SavedImageResponse{
			Status: false,
			ErrorMsg: fetchUserErr.Error(),
		}
	}

	if fetchedUser.ProfilePicture != "" {
		if err := removeProfilePictureForUser(ctx, fetchedUser, dao, redisUserOp); err != nil {
			return &models.SavedImageResponse{
				Status: false,
				ErrorMsg: err.Error(),
			}
		}
	}

	awsUUID := uuid.NewString()

	fetchedUser.ProfilePicture = awsUUID

	fileObject, err := saveImageAWSRedisDB(ctx, fetchedUser, imageFile, imageDet, dao, redisUserOp); 
	if err != nil {
		return &models.SavedImageResponse{
			Status: false,
			ErrorMsg: err.Error(),
		}
	}

	return &models.SavedImageResponse{
		Status: true,
		FileId: fileObject.Id,
		ImageURL: fileObject.SecureURL,
	}
}