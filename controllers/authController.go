package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/leksyking/go-authentication/models"
	"github.com/leksyking/go-authentication/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	validate                         = validator.New()
	Client         *mongo.Client     = models.Client
	UserCollection *mongo.Collection = models.UserCollection(Client)
)

func hashPassword(userPassword string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(userPassword), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
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
	exists, err := UserCollection.CountDocuments(ctx, bson.M{"email": User.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	if exists > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email exists already"})
		return
	}
	//hash passord
	hash := hashPassword(*User.Password)
	User.Password = &hash
	//save user
	_, err = UserCollection.InsertOne(ctx, User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	userId := User.ID.Hex()
	//create token(with id)
	//attach cookies to user(use jwt to create the tokens to be stored in cookies)
	accessToken, refreshToken, _ := utils.GenerateToken(*User.Email, *User.UserName, userId)
	utils.AttachCookiesToResponse(accessToken, refreshToken, c)
	//send verification token to user's email
	c.JSON(http.StatusCreated, gin.H{"msg": "Successful..."})
}

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Login Users"})
}
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Logout Users"})
}
