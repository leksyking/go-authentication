package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leksyking/go-authentication/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Client         *mongo.Client     = models.Client
	UserCollection *mongo.Collection = models.UserCollection(Client)
)

func hashPassword() {
}
func verifyPassword() {

}

func Register(c *gin.Context) {
	//check if email exists already
	
	//validate strut
	//hash passord
	//save user
	//create token(with id, and username)
	//attach cookies to user(use jwt to use te tokem to be stored in cookies)
	//send verification token to user's email
	c.JSON(http.StatusCreated, gin.H{"msg": "Register Users"})
}
func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Login Users"})
}
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Logout Users"})
}
