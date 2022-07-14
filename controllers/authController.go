package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base32"
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

func verifyPassword(hashedPassword, enteredPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(enteredPassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Incorrect Password"
		valid = false
	}
	return valid, msg
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
	//generate verification token
	randomBytes := make([]byte, 50)
	_, err = rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	verificationToken := base32.StdEncoding.EncodeToString(randomBytes)[:40]
	User.VerificationToken = &verificationToken
	//save user
	_, insertErr := UserCollection.InsertOne(ctx, User)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
		fmt.Println(err)
		return
	}
	userId := User.ID.Hex()
	//attach cookies to user(use jwt to create the tokens to be stored in cookies)
	accessToken, refreshToken, _ := utils.GenerateToken(*User.Email, *User.UserName, userId)
	utils.AttachCookiesToResponse(accessToken, refreshToken, c)
	//send verification token to user's email
	origin := "http://localhost:8080/api/v1"
	email := []string{*User.Email}
	err = utils.SendVerificationEmail(origin, verificationToken, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"msg": "Successful..., check your mail to verify your account"})
}

//verifyemail
//check for the token and change status
//forgotPassword
//get the email
//send forgotPassword link
//resetPassword
//change the password
func Login(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	//check for the email in the database
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	var foundUser models.User
	err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}
	valid, msg := verifyPassword(*foundUser.Password, *user.Password)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
	}
	//attach cookies
	userId := foundUser.ID.Hex()
	accessToken, refreshToken, _ := utils.GenerateToken(*foundUser.Email, *foundUser.UserName, userId)
	utils.AttachCookiesToResponse(accessToken, refreshToken, c)
	c.JSON(http.StatusOK, gin.H{"msg": "Login Users"})
}
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Logout Users"})
}
