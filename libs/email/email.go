package libs

import (
	"errors"
	"fmt"
	"os"
	"profiln-be/libs"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type IEmail interface {
	SendAuthEmail(subject string, to []string, templateData any, templateFilename string) error
}

type Email struct {
	SenderEmail  string
	SMTPHost     string
	SMTPPort     int
	AuthEmail    string
	AuthPassword string

	log *logrus.Logger
}

func NewEmail(port int, sender, host, email, password string, log *logrus.Logger) IEmail {
	return &Email{
		SenderEmail:  sender,
		SMTPHost:     host,
		SMTPPort:     port,
		AuthEmail:    email,
		AuthPassword: password,
		log:          log,
	}
}

func (e *Email) SendAuthEmail(subject string, to []string, templateData any, templateFilename string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		errMsg := fmt.Sprintf("os.Getwd: %v", err)
		err = errors.New(errMsg)

		return err
	}

	filepath := fmt.Sprintf("%s/libs/email/template/%s", workingDir, templateFilename)

	body, err := libs.HTMLToString(filepath, templateData)
	if err != nil {
		errMsg := fmt.Sprintf("libs.HTMLToString: %v", err)
		err = errors.New(errMsg)

		return err
	}

	err = e.sendEmail(subject, to, body)
	if err != nil {
		errMsg := fmt.Sprintf("sendEmail: %v", err)
		err = errors.New(errMsg)

		return err
	}

	e.log.Infof("success send '%s' email", subject)
	return nil
}

func (e *Email) sendEmail(subject string, to []string, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", e.SenderEmail)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		e.SMTPHost,
		e.SMTPPort,
		e.AuthEmail,
		e.AuthPassword,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	e.log.Info("success send email")
	return nil
}
