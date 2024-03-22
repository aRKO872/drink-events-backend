package pkg_helpers

import (
	"fmt"
	"net/smtp"
	"regexp"

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
