package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading environment variables")
	}
	server := gin.New()
	server.Use(gin.Logger(), gin.Recovery())

	port := os.Getenv("PORT")
	server.Run(":" + port)
}
