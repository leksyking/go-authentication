package utils

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/leksyking/go-authentication/wait"
)

var (
	username = os.Getenv("SMTP_USERNAME")
	password = os.Getenv("SMTP_PASSWORD")
	host     = os.Getenv("SMTP_HOST")
	wg       = wait.Wg
)

func PlainAuth() smtp.Auth {
	auth := smtp.PlainAuth("GO-AUTH-LEKSYBABA", username, password, host)
	return auth
}

func SendVerificationEmail(origin, verificationToken string, email []string, c *gin.Context) {
	defer wg.Done()
	a := PlainAuth()
	addr := host + ":587"
	verifyEmail := fmt.Sprintf("%s/auth/verify-email?token=%s&email=%s", origin, verificationToken, email)
	subject := fmt.Sprintf("Please confirm your email by clicking this link: %s", verifyEmail)

	message := []byte("To: " + email[0] + "\r\n" +
		"Subject: Welcome to Go website!\r\n" +
		"\r\n" +
		subject + ".\r\n")

	err := smtp.SendMail(addr, a, "Leksyking", email, message)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again later..."})
	}
}
