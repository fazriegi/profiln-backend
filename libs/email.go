package libs

import (
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(subject string, to []string, body string) error {
	senderMail := os.Getenv("SENDER_MAIL")
	const CONFIG_SMTP_HOST = "smtp.gmail.com"
	const CONFIG_SMTP_PORT = 587
	CONFIG_AUTH_EMAIL := senderMail
	CONFIG_AUTH_PASSWORD := os.Getenv("SENDER_PASSWORD")

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
