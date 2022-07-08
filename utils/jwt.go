package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

//create token
type SignedDetails struct {
	Email    string
	UserName string
	Uid      string
	jwt.StandardClaims
}

var (
	SECRET_KEY = os.Getenv("JWT_SECRET")
)

func GenerateToken(email, username, id string) (string, string, error) {
	accessClaims := &SignedDetails{
		Email:    email,
		UserName: username,
		Uid:      id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		//add user details to refresh and generate refresh token
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(720)).Unix(),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

//validate token
func AttachCookiesToResponse() {
	//create access jwt token
	//create refresh jwt tokenpackage main

	// import (
	// 	"github.com/gin-contrib/sessions"
	// 	"github.com/gin-contrib/sessions/cookie"
	// 	"github.com/gin-gonic/gin"
	// 	)

	// 	func main() {
	// 		r := gin.Default()
	// 		store := cookie.NewStore([]byte("secret"))
	// 		store.Options(sessions.Options{MaxAge:   60 * 60 * 24}) // expire in a day
	// 		r.Use(sessions.Sessions("mysession", store))

	// 		r.GET("/incr", func(c *gin.Context) {
	// 			session := sessions.Default(c)
	// 			var count int
	// 			v := session.Get("count")
	// 			if v == nil {
	// 				count = 0
	// 			} else {
	// 	// 				count = v.(int)
	// 	// 				count++
	// 	// 			}
	// 	// 			session.Set("count", count)
	// 	// 			session.Save()
	// 	// 			c.JSON(200, gin.H{"count": count})
	// 	// 		})
	// 	// 		r.Run(":8000")
	// 	// 	}
	// 	outer := gin.Default();

	// token_value := func(c *gin.Context) string {

	//     var value string

	//     if cookie, err := c.Request.Cookie("session"); err == nil {
	//       value = cookie.Value
	//     } else {
	//       value = RandToken(64)
	//     }
	//     return value
	//   }

	//   cookie_store := cookie.NewStore([]byte(token_value))
	//   router.Use(sessions.Sessions("session",cookie_store))
	//attach cookies
}
