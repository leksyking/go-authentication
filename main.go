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
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading environment variables")
	}
	server := gin.New()
	server.Use(gin.Logger(), gin.Recovery())

	//routes
	routes.AuthRouter(server)
	routes.UserRouter(server)
	server.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Welcome to Go Authentication")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server.Run(":" + port)
}
