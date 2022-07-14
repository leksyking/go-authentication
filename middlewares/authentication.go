package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
	cookies, _ := c.Cookie("accessCookie")
	fmt.Println(cookies)
	// accessToken := cookies[0]
	// refreshToken := cookies[1]
	//check for cookies
	//accesstoken first
	//refreshtoken next
	//verify the token inside  the cookies
	//store the payloads in user
	//next
}
