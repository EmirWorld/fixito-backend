package config

import "fixito-backend/docs"

func SetupSwagger() {
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "FIXITO BACKED API"
	docs.SwaggerInfo.Description = "This is a server for Fixito"
	docs.SwaggerInfo.Version = "1.0"
	//docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

}
