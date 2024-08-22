package pkg_helpers

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"regexp"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	pkg_config "github.com/drink-events-backend/pkg/config"
)

func IsValidEmail(email string) bool {
	// Regular expression to match email pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidPhoneNumber(phoneNumber string) bool {
	// Regular expression to match Indian phone numbers starting with 6-9
	phoneRegex := regexp.MustCompile(`^[6-9]\d{9}$`)
	return phoneRegex.MatchString(phoneNumber)
}

func SendEmail(to []string, mailBody []byte) error {
	configInfo := pkg_config.GetProjectConfig()
	auth := smtp.PlainAuth(
		"",
		configInfo.SMTP_USER,
		configInfo.SMTP_PASSWORD,
		configInfo.SMTP_HOST,
	)

	if sendEmailErr := smtp.SendMail(
		fmt.Sprintf("%s:%d", configInfo.SMTP_HOST, configInfo.SMTP_PORT),
		auth,
		configInfo.SMTP_USER,
		to,
		mailBody,
	); sendEmailErr != nil {
		fmt.Println("Error sending mail : ", sendEmailErr)
		return sendEmailErr
	}

	return nil
}

func GetUserDetailsFromReqHeader(
	req *http.Request,
) (*models.TokenizedUserDetails, error) {
	userId := req.Header.Get(literals.HEADER_USER_ID)
	userType := req.Header.Get(literals.HEADER_USER_TYPE)
	radius := req.Header.Get(literals.HEADER_USER_SEARCH_DISTANCE)

	if userId == "" || userType == "" || radius == "" {
		return &models.TokenizedUserDetails{}, fmt.Errorf("unable to resolve user dtls from header")
	}

	return &models.TokenizedUserDetails{
		UserType: userType,
		UserId: userId,
	}, nil
}

func ValidateImageRequest(req *http.Request) (multipart.File, *multipart.FileHeader, error) {
	err := req.ParseMultipartForm(120 << 10) // 120 kb limit
	if err != nil {
		return nil, nil, errors.New("unable to parse form")
	}

	// Get the file from the form
	file, fileHandler, err := req.FormFile("picture")
	if err != nil {
		return nil, nil, errors.New("error retrieving the file")
	}
	defer file.Close()

	// Validate file size (max 100 kb)
	if fileHandler.Size > 100 << 10 {
		return nil, nil, errors.New("file too large. Max size is 5MB")
	}

	// Validate file type (allow only jpeg, jpg and png)
	fileType := fileHandler.Header.Get("Content-Type")
	if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/jpg" {
		return nil, nil, errors.New("invalid file type. Only JPEG, JPG and PNG are allowed")
	}

	return file, fileHandler, nil
}

func Smaller(id1, id2 string) (smallId, bigId string) {
	if id1 < id2 {
		smallId = id1
		bigId = id2
	} else {
		smallId = id2
		bigId = id1
	}

	return
}