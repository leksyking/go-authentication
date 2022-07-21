package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leksyking/go-authentication/controllers"
	"github.com/leksyking/go-authentication/middlewares"
)

func AuthRouter(route *gin.Engine) {
	authRoute := route.Group("/api/v1/auth")
	{
		authRoute.POST("/register", controllers.Register)
		authRoute.POST("/login", controllers.Login)
		authRoute.POST("/verify-email", controllers.VerifyEmail)
		authRoute.POST("/logout", middlewares.Authentication, controllers.Logout)
		authRoute.POST("/forgot-password", controllers.ForgotPassword)
		authRoute.POST("/reset-password", controllers.ResetPassword)
	}
}
