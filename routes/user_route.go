package routes

import (
	"github.com/etg-dev/restApi/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	userGroups := router.Group("/api/users")
	{
		userGroups.GET("/:id", controllers.GetUser())
		userGroups.PUT("/:id", controllers.UpdateUser())
		userGroups.DELETE("/:id", controllers.DeleteUser())
		userGroups.GET("/", controllers.GetUsers())
		userGroups.POST("/", controllers.CreateUser())
	}
}
