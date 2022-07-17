package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leksyking/go-authentication/utils"
)

func Authentication(c *gin.Context) {
	accessToken, err := c.Cookie("accessCookie")
	if err != nil {
		refreshToken, err := c.Cookie("refreshCookie")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed."})
		}
		payload, msg := utils.ValidateToken(refreshToken)
		if msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		}
	}
	payload, msg := utils.ValidateToken(accessToken)
	if msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
	}
	//c.Request.Response = payload
	c.Next()
	//store the payloads in user
	//next
}
