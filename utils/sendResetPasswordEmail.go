package utils

import (
	"fmt"
	"net/http"
	"net/smtp"
	"sync"

	"github.com/gin-gonic/gin"
)

var wg *sync.WaitGroup

func SendResetPasswordEmail(origin, passwordToken string, email []string, c *gin.Context) {
	a := PlainAuth()
	addr := host + ":587"
	resetPassword := fmt.Sprintf("%s/auth/reset-password?email=%s&token=%s", origin, email, passwordToken)
	subject := fmt.Sprintf("<p>Please reset your passsword by clicking on this link <a href=%s>Reset Password</a></p>", resetPassword)

	message := []byte("To: " + email[0] + "\r\n" +
		"Subject: Forgot Password!\r\n" +
		"\r\n" +
		subject + ".\r\n")
	err := smtp.SendMail(addr, a, "Olukoya", email, message)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again later..."})
	}
	wg.Done()
}
