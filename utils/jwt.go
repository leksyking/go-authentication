package utils

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

//create token
type SignedDetails struct {
	Email    string
	UserName string
	Uid      string
	jwt.RegisteredClaims
}

var (
	SECRET_KEY = os.Getenv("JWT_SECRET")
)

func GenerateToken(email, username, id string) (string, string, error) {
	accessClaims := &SignedDetails{
		Email:    email,
		UserName: username,
		Uid:      id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(24)).Local()),
		},
	}
	refreshClaims := &SignedDetails{
		//add user details to refresh and generate refresh token
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(720)).Local()),
		},
	}
	accessTokenJWT, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return "", "", err
	}
	refreshTokenJWT, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return "", "", err
	}
	return accessTokenJWT, refreshTokenJWT, nil
}

//validate token

func AttachCookiesToResponse(accessTokenJWT, refreshTokenJWT string, c *gin.Context) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	c.SetCookie("accessCookie", accessTokenJWT, 60*60*24, "/", "localhost", env == "development", true)
	c.SetCookie("refreshCookie", refreshTokenJWT, 60*60*24*30, "/", "localhost", env == "development", true)
}
