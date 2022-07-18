package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leksyking/go-authentication/utils"
)

func ShowUser(c *gin.Context) {
	payload, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed."})
	}
	Uid := payload.(*utils.SignedDetails).Uid
	Username := payload.(*utils.SignedDetails).UserName
	details := fmt.Sprintf("%s,\n %s", Uid, Username)
	c.JSON(http.StatusOK, gin.H{"user": details})
}
