package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leksyking/go-authentication/utils"
)

func ShowUser(c *gin.Context) {
	payload, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed."})
	}
	Username := payload.(*utils.SignedDetails).UserName
	c.JSON(http.StatusOK, gin.H{"user": Username})
}
