package main

import (
	"github.com/Bahnstar/petitionu-api/controllers"
	"github.com/Bahnstar/petitionu-api/initializers"
	"github.com/Bahnstar/petitionu-api/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectDatabase()
	initializers.SyncDb()
}

func main() {
	r := gin.Default()

	r.POST("/sign-up", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.Run()
}
