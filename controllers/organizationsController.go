package controllers

import (
	"net/http"

	"github.com/Bahnstar/petitionu-api/initializers"
	"github.com/Bahnstar/petitionu-api/models"
	"github.com/gin-gonic/gin"
)

func CreateOrgraniztion(c *gin.Context) {
	var body struct {
		Name        string
		Description string
		Users       []models.User
		Petitions   []models.Petition
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	organization := models.Organization{Name: body.Name, Description: body.Description, Users: body.Users, Petitions: body.Petitions}

	result := initializers.DB.Create(&organization)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization created successfully"})
}

func GetOrganizations(c *gin.Context) {
	var organizations []models.Organization

	if result := initializers.DB.Find(&organizations); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
	}

	c.JSON(http.StatusOK, organizations)
}

func GetOrganization(c *gin.Context) {
	id := c.Param("id")

	var organization models.Organization

	if result := initializers.DB.First(&organization, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, organization)
}
