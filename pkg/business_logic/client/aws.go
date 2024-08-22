package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/drink-events-backend/models"
)

func SaveObjectToAWS(
	awsUUID string,
	imageFile multipart.File,
	imageDet *multipart.FileHeader,
) (models.FileObject, error) {
	awsRegion := "eu-north-1"
	bucketName := "drink-events-images"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		return models.FileObject{}, errors.New("error setting up session")
	}

	awsServe := s3.New(sess)
	fileBuffer := new(bytes.Buffer)

	if _, err := io.Copy(fileBuffer, imageFile); err != nil {
		return models.FileObject{}, errors.New("error reading image file")
	}

	awsOutput, awsErr := awsServe.PutObject(&s3.PutObjectInput{
		Body: bytes.NewReader(fileBuffer.Bytes()),
		Key: aws.String(awsUUID),
		Bucket: aws.String(bucketName),
	})

	secureURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, awsRegion, awsUUID)

	if awsErr != nil {
		return models.FileObject{}, errors.New("error adding file to aws")
	}

	return models.FileObject{
		Etag: *awsOutput.ETag,
		FileName: imageDet.Filename,
		AWS_Key: awsUUID,
		Id: awsUUID,
		SecureURL: secureURL,
	}, nil
}

func DeleteObjectFromAWS(awsUUID string) error {
	awsRegion := "eu-north-1"
	bucketName := "drink-events-images"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		return errors.New("error setting up session")
	}

	awsServe := s3.New(sess)

	if _, awsErr := awsServe.DeleteObject(&s3.DeleteObjectInput{
		Key: aws.String(awsUUID),
		Bucket: aws.String(bucketName),
	}); awsErr != nil {
		return errors.New("error deleting file from aws")
	}

	return nil
}