package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leksyking/go-authentication/controllers"
	"github.com/leksyking/go-authentication/middlewares"
)

func UserRouter(route *gin.Engine) {
	userRoute := route.Group("/api/v1/user", middlewares.Authentication)
	{
		userRoute.GET("/show", controllers.ShowUser)
	}
}
