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
	router := gin.Default()

	router.POST("/sign-up", controllers.SignUp)
	router.POST("/login", controllers.Login)
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)

	users := router.Group("/users")
	{
		users.GET("", controllers.GetUsers)
		users.GET("/:id", controllers.GetUser)
		users.PATCH("/:id", controllers.UpdateUser)
		users.DELETE("/:id", controllers.DeleteUser)
	}

	organizations := router.Group("/organizations")
	{
		organizations.GET("", controllers.GetOrganizations)
		organizations.GET("/:id", controllers.GetOrganization)
		organizations.POST("", controllers.CreateOrganization)
		organizations.PATCH("/:id", controllers.UpdateOrganization)
		organizations.DELETE("/:id", controllers.DeleteOrganization)
	}

	preferences := router.Group("/preferences")
	{
		preferences.GET("", controllers.GetPreferences)
		preferences.GET("/:id", controllers.GetPreference)
		preferences.POST("", controllers.CreatePreference)
		preferences.PATCH("/:id", controllers.UpdatePreference)
		preferences.DELETE("/:id", controllers.DeletePreference)
	}

	petitions := router.Group("/petitions")
	{
		petitions.GET("", controllers.GetPetitions)
		petitions.GET("/:id", controllers.GetPetition)
		petitions.POST("", controllers.CreatePetition)
		petitions.PATCH("/:id", controllers.UpdatePetition)
		petitions.DELETE("/:id", controllers.DeletePetition)
	}

	comments := router.Group("/comments")
	{
		comments.GET("", controllers.GetComments)
		comments.GET("/:id", controllers.GetComment)
		comments.POST("", controllers.CreateComment)
		comments.PATCH("/:id", controllers.UpdateComment)
		comments.DELETE("/:id", controllers.DeleteComment)
	}

	router.Run()
}
