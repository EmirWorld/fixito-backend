package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"poosible-backend/controllers"
	"poosible-backend/middleware"
)

// SetupRouter is a function to set up all routes
var SetupRouter = func(router *gin.Engine) {
	v1api := router.Group("/v1/api")
	{

		//PUBLIC ROUTER

		//Auth Routes
		v1api.POST("auth/login", controllers.Login())
		v1api.POST("auth/logout", controllers.Logout())

		//User Routes
		v1api.POST("/user", controllers.CreateUser())

		//Helper Routes
		v1api.GET("/helper/countries", controllers.Countries())
		v1api.GET("/helper/currencies", controllers.Currencies())

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
			protectedGroup.GET("/organisation", controllers.GetOrganisation())
			protectedGroup.PUT("/organisation/:organisationId", controllers.UpdateOrganisation())
			protectedGroup.DELETE("/organisation/:organisationId", controllers.DeleteOrganisation())

			//Item Routes
			protectedGroup.POST("/item", controllers.CreateItem())
			protectedGroup.GET("/item/:itemId", controllers.GetItem())
			protectedGroup.GET("/items", controllers.GetItems())
			protectedGroup.PUT("/item/:itemId", controllers.UpdateItem())
			protectedGroup.DELETE("/item/:itemId", controllers.DeleteItem())
		}

	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
