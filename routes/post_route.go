package routes

import (
	"github.com/etg-dev/restApi/controllers"
	"github.com/etg-dev/restApi/middleware"
	"github.com/gin-gonic/gin"
)

func PostRoute(router *gin.Engine) {
	postGroup := router.Group("/api/posts")
	{
		postGroup.GET("/user/:userId", controllers.GetUsersPosts())
		postGroup.GET("/:userId", middleware.ValidateUserID(), middleware.ValidateAction("Read"), middleware.Paginate(), controllers.GetPosts())
		postGroup.POST("/:userId", controllers.CreatePost())
	}

}
