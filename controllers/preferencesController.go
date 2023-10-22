package controllers

import (
	"net/http"

	"github.com/Bahnstar/petitionu-api/initializers"
	"github.com/Bahnstar/petitionu-api/models"
	"github.com/gin-gonic/gin"
)

type CreatePreferenceBody struct {
	Name   string
	Value  string
	UserId uint
}

func GetPreferences(c *gin.Context) {
	var preferences []models.Preference

	if result := initializers.DB.Find(&preferences); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
	}

	c.JSON(http.StatusOK, gin.H{"data": &preferences})
}

func GetPreference(c *gin.Context) {
	id := c.Param("id")
	var preference models.Preference

	if result := initializers.DB.First(&preference, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": &preference})
}

func CreatePreference(c *gin.Context) {
	var body CreatePreferenceBody

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	preference := models.Preference{Name: body.Name, Value: body.Value, UserId: body.UserId}

	result := initializers.DB.Create(&preference)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create preference"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Preference created successfully"})
}
