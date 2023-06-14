package router

import (
	"fixito-backend/controllers"
	"fixito-backend/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var SetupRouter = func(router *gin.Engine) {
	v1api := router.Group("/v1/api")
	{

		//PUBLIC ROUTER

		//Auth Routes
		v1api.POST("auth/login", controllers.Login())
		v1api.POST("auth/logout", controllers.Logout())

		//User Routes
		v1api.POST("/user", controllers.CreateUser())

		// Apply the AuthMiddleware to the protected group
		protectedGroup := v1api.Group("", middleware.AuthMiddleware())
		{
			// PRIVATE ROUTER

			//User Routes
			protectedGroup.GET("/user/current", controllers.GetCurrentUser())
			protectedGroup.GET("/user/:userId", controllers.GetUser())
			protectedGroup.PUT("/user/:userId", controllers.UpdateUser())
			protectedGroup.DELETE("/user/:userId", controllers.DeleteUser())

			//Organization Routes
			protectedGroup.POST("/organisation", controllers.CreateOrganisation())
			protectedGroup.GET("/organisation/:organisationId", controllers.GetOrganisation())
			protectedGroup.PUT("/organisation/:organisationId", controllers.UpdateOrganisation())

		}

	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
