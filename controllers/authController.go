package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/leksyking/go-authentication/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	validate                         = validator.New()
	Client         *mongo.Client     = models.Client
	UserCollection *mongo.Collection = models.UserCollection(Client)
)

func hashPassword() {
	//bcrypt.GenerateFromPassword()
}
func verifyPassword() {
	//bcrypt.CompareHashAndPassword()
}

func Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var User models.User
	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	User.ID = primitive.NewObjectID()
	//validate struct
	err := validate.Struct(User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	//check if email exists already
	var emailExists models.Email
	err = UserCollection.FindOne(ctx, bson.M{"email": User.Email}).Decode(&emailExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	if emailExists.Email != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email exists already"})
		return
	}

	//hash passord
	//save user
	_, err = UserCollection.InsertOne(ctx, User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	//create token(with id, and username)
	//attach cookies to user(use jwt to use te tokem to be stored in cookies)
	//send verification token to user's email
	c.JSON(http.StatusCreated, gin.H{"msg": "Successful..."})
}
func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Login Users"})
}
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Logout Users"})
}
