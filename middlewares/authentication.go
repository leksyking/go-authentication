package middlewares

import (
	"fmt"
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
		//attach cookies again and assign payload
		c.Set("user", payload)
		c.Next()
	}
	// c.Set("user",) & c.get("user")   uid := user.(*model.User).UID
	payload, msg := utils.ValidateToken(accessToken)
	if msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
	}
	fmt.Println(payload)
	c.Set("user", payload)
	//store the payloads in user
	//c.Request.Response = payload
	c.Next()
}
