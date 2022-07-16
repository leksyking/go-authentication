package middlewares

import (
	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
	accessToken, err := c.Cookie("accessCookie")

	refreshToken, err := c.Cookie("refreshCookie")
	// accessToken := cookies[0]
	// refreshToken := cookies[1]
	//check for cookies
	//accesstoken first
	//refreshtoken next
	//verify the token inside  the cookies
	//store the payloads in user
	//next
}
