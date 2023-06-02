package router

import (
	"fixito-backend/controllers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var SetupRouter = func(router *gin.Engine) {
	v1api := router.Group("/v1/api")
	{
		//Auth Routes
		v1api.POST("auth/login", controllers.Login())

		//User Routes
		v1api.POST("/user", controllers.CreateUser())
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
