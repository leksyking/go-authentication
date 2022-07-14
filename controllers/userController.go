package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Show user"})
}
