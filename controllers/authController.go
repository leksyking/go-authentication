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

}
func Login(c *gin.Context) {

}
func Logout(c *gin.Context) {

}
