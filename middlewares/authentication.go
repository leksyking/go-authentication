package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leksyking/go-authentication/models"
	"github.com/leksyking/go-authentication/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Client          *mongo.Client     = models.Client
	TokenCollection *mongo.Collection = models.TokenCollection(Client)
)

func Authentication(c *gin.Context) {
	accessToken, err := c.Cookie("accessCookie")
	if err != nil {
		refreshToken, err := c.Cookie("refreshCookie")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		}
		payload, msg := utils.ValidateToken(refreshToken)
		if msg != "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		}
		var token models.Token
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userId := payload.Uid
		usertId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops!, Something went wrong."})
		}
		refresht_Token := payload.RefreshToken
		err = TokenCollection.FindOne(ctx, bson.M{"user_id": usertId, "refreshtoken": refresht_Token}).Decode(&token)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops!, Something went wrong."})
		}
		accessTokenJWT, refreshTokenJWT, _ := utils.GenerateToken(payload.Email, payload.UserName, refreshToken, userId)
		utils.AttachCookiesToResponse(accessTokenJWT, refreshTokenJWT, c)
		c.Set("user", payload)
		c.Next()
	}
	payload, msg := utils.ValidateToken(accessToken)
	if msg != "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
	}
	c.Set("user", payload)
	c.Next()
}
