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
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(signedToken *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		msg = err.Error()
		return nil, msg
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "Token is invalid"
		return nil, msg
	}
	return claims, ""
}

func AttachCookiesToResponse(accessTokenJWT, refreshTokenJWT string, c *gin.Context) {
	// secure is set to false if development is local
	acookie := accessTokenJWT
	c.SetCookie("accessCookie", acookie, 60*60*24, "/", "localhost", false, true)

	rcookie := refreshTokenJWT
	c.SetCookie("refreshCookie", rcookie, 60*60*24*30, "/", "localhost", false, true)
}
