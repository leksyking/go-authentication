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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		}
		payload, msg := utils.ValidateToken(refreshToken)
		if msg != "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		}
		//attach cookies again and assign payload
		c.Set("user", payload)
		c.Next()
	}
	payload, msg := utils.ValidateToken(accessToken)
	if msg != "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
	}
	c.Set("user", payload)
	c.Next()
}
