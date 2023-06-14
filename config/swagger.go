package config

import "fixito-backend/docs"

func SetupSwagger() {
	// Set Swagger information programmatically
	swaggerTitle := "FIXITO API"
	swaggerDescription := "FIXITO is a compact and efficient Point of Sale (POS) system specifically designed for small businesses, particularly those operating in the automotive repair industry, such as car fix shops. This application serves as a comprehensive solution that helps streamline business operations and enhance efficiency for these establishments."
	swaggerVersion := "1.0"
	swaggerBasePath := "/v1"
	swaggerSchemes := []string{"http", "https"}

	// Update SwaggerInfo properties
	docs.SwaggerInfo.Title = swaggerTitle
	docs.SwaggerInfo.Description = swaggerDescription
	docs.SwaggerInfo.Version = swaggerVersion
	docs.SwaggerInfo.BasePath = swaggerBasePath
	docs.SwaggerInfo.Schemes = swaggerSchemes
}
