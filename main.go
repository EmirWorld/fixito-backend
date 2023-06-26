package main

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"poosible-backend/config"
	"poosible-backend/router"

	"time"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	r := gin.New()
	gob.Register(time.Time{})
	authStore := cookie.NewStore([]byte("auth-secret"))
	r.Use(sessions.Sessions("auth-session", authStore))
	config.ConnectDatabase()
	config.SetupSwagger()
	router.SetupRouter(r)
	r.Run("0.0.0.0:9090")
}
