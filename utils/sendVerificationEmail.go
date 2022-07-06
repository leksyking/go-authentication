package utils

import "net/smtp"

func PlainAuth() smtp.Auth {
	auth := smtp.PlainAuth()
	return auth
}

func SendEmail() error {
	err := smtp.SendMail()
	return err
}
