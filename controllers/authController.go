package controllers

import (
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
	//validate struct
	//hash password
	//save user
	//create token(with id, and username)
	//attach cookies to user(use jwt to use the tokem to be stored in cookies)
	//send verification token to user's email
}
func Login(c *gin.Context) {

}
func Logout(c *gin.Context) {

}
