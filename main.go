package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.New()
	server.Use(gin.Logger(), gin.Recovery())

	port := os.Getenv("PORT")
	server.Run(":" + port)
}
