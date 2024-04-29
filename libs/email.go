package libs

import (
	"os"
	"strconv"
	"gopkg.in/gomail.v2"
)

func SendEmail(subject string, to []string, body string) error {
	senderMail := os.Getenv("AUTH_EMAIL")
	CONFIG_SMTP_HOST := os.Getenv("SMTP_HOST")
	CONFIG_SMTP_PORT,_ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	CONFIG_AUTH_EMAIL := os.Getenv("AUTH_EMAIL")
	CONFIG_AUTH_PASSWORD := os.Getenv("AUTH_PASSWORD")

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", senderMail)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}
	return nil
}
