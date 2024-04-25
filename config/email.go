package config

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

var MailDialer *gomail.Dialer

func NewMailDialer() {
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	account := os.Getenv("AUTH_EMAIL")
	password := os.Getenv("AUTH_PASSWORD")

	MailDialer = gomail.NewDialer(
		host,
		port,
		account,
		password,
	)
}
