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
	validate                          = validator.New()
	Client          *mongo.Client     = models.Client
	UserCollection  *mongo.Collection = models.UserCollection(Client)
	TokenCollection *mongo.Collection = models.TokenCollection(Client)
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

//verify token
func VerifyEmail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	var foundUser models.User
	err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	if *foundUser.VerificationToken != *user.VerificationToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		fmt.Println("Invalid token")
		return
	}
	time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	_, err = UserCollection.UpdateOne(ctx, bson.D{primitive.E{Key: "_id", Value: foundUser.ID}},
		bson.D{{Key: "$set", Value: bson.M{"verification_token": "", "is_verified": true, "verified": time}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Email verified"})
}

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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		return
	}
	valid, msg := verifyPassword(*foundUser.Password, *user.Password)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
		return
	}
	//check whether user is verified
	if !*foundUser.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid verification"})
		fmt.Println("invalid verification")
		return
	}
	userId := foundUser.ID.Hex()

	//check for user in the token collection
	refreshToken := ""
	randomBytes := make([]byte, 50)
	_, err = rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	token, err := TokenCollection.CountDocuments(ctx, bson.D{primitive.E{Key: "user_id", Value: foundUser.ID}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	fmt.Println(token)
	if token > 0 {
		var tokenUser models.Token
		if err := TokenCollection.FindOne(ctx, bson.M{"user_id": foundUser.ID}).Decode(&tokenUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't decode token."})
			return
		}
		if !tokenUser.IsValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid details"})
			return
		}
		refreshToken = *tokenUser.RefreshToken
		accessTokenJWT, refreshTokenJWT, _ := utils.GenerateToken(*foundUser.Email, *foundUser.UserName, refreshToken, userId)
		utils.AttachCookiesToResponse(accessTokenJWT, refreshTokenJWT, c)
		c.JSON(http.StatusOK, gin.H{"msg": "Login Successful"})
		return
	}
	var userToken models.Token
	refreshToken = base32.StdEncoding.EncodeToString(randomBytes)[:40]
	userAgent := c.Request.Header["User-Agent"][0]
	ip := c.ClientIP()
	userToken.TokenID = primitive.NewObjectID()
	userToken.RefreshToken = &refreshToken
	userToken.UserAgent = &userAgent
	userToken.IP = &ip
	userToken.IsValid = true
	userToken.User = foundUser.ID
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = TokenCollection.InsertOne(ctx, userToken)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving token."})
		return
	}
	accessTokenJWT, refreshTokenJWT, _ := utils.GenerateToken(*foundUser.Email, *foundUser.UserName, refreshToken, userId)
	utils.AttachCookiesToResponse(accessTokenJWT, refreshTokenJWT, c)
	c.JSON(http.StatusOK, gin.H{"msg": "Login Successful"})
}

func Logout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed."})
		return
	}
	userId := user.(*utils.SignedDetails).ID
	usertId, _ := primitive.ObjectIDFromHex(userId)
	_, err := TokenCollection.DeleteOne(ctx, bson.M{"user_id": usertId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again later."})
		fmt.Println(err)
		return
	}
	c.SetCookie("accessCookie", "logout", 0, "/", "localhost", false, true)
	c.SetCookie("refreshCookie", "logout", 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"msg": "You are logged out"})
}

func ForgotPassword(c *gin.Context) {
	//get email
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again later."})
		fmt.Println(err)
		return
	}
	//verify whether email is valid in the db
	err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
		fmt.Println(err)
		return
	}
	//send resetpassword mail

}

func ResetPassword(c *gin.Context) {

}
