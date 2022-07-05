package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leksyking/go-authentication/controllers"
)

func AuthRouter(route *gin.Engine) {
	authRoute := route.GROUP("/api/v1/auth")
	{
		authRoute.POST("/register", controllers.Register)
		authRoute.POST("/login", controllers.Login)
		authRoute.POST("/logout", controllers.Logout)
	}
}
