package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/leksyking/go-authentication/routes"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading environment variables")
	}
	server := gin.New()
	server.Use(gin.Logger(), gin.Recovery())

	//routes
	server.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Welcome to Go Authentication")
	})
	routes.AuthRouter(server)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server.Run(":" + port)
}
