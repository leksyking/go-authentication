package utils

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	username = os.Getenv("SMTP_USERNAME")
	password = os.Getenv("SMTP_PASSWORD")
	host     = os.Getenv("SMTP_HOST")
)

func PlainAuth() smtp.Auth {
	auth := smtp.PlainAuth("GO-AUTH-LEKSYBABA", username, password, host)
	return auth
}

func SendVerificationEmail(origin, verificationToken string, email []string, c *gin.Context) {
	a := PlainAuth()
	addr := host + ":587"
	verifyEmail := fmt.Sprintf("%s/auth/verify-email?token=%s&email=%s", origin, verificationToken, email)
	subject := fmt.Sprintf("<p>Please confirm your email by clicking this link: <a href=%s>Verify Email</a></p>", verifyEmail)

	message := []byte("To: " + email[0] + "\r\n" +
		"Subject: Welcome to Go website!\r\n" +
		"\r\n" +
		subject + ".\r\n")

	err := smtp.SendMail(addr, a, "Leksyking", email, message)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again later..."})
	}
	wg.Done()
}
