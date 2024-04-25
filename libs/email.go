package libs

import (
	"os"
	"profiln-be/config"

	"gopkg.in/gomail.v2"
)

func SendMail(subject string, to []string, body string) error {
	sender := os.Getenv("SENDER_NAME")

	message := gomail.NewMessage()
	message.SetHeader("From", sender)
	message.SetHeader("To", to...)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	err := config.MailDialer.DialAndSend(message)
	if err != nil {
		return err
	}

	return nil
}
