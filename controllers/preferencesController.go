package controllers

import (
	"net/http"

	"github.com/Bahnstar/petitionu-api/initializers"
	"github.com/Bahnstar/petitionu-api/models"
	"github.com/gin-gonic/gin"
)

type CreatePreferenceBody struct {
	Name   string `binding:"required"`
	Value  string
	UserId uint `binding:"required"`
}

type UpdatePreferenceBody struct {
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

	var user models.User

	if result := initializers.DB.First(&user, body.UserId); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	preference := models.Preference{Name: body.Name, Value: body.Value, UserId: user.ID}

	result := initializers.DB.Create(&preference)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create preference"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Preference created successfully"})
}

func UpdatePreference(c *gin.Context) {
	id := c.Param("id")
	var body UpdatePreferenceBody

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body in update preference"})
	}

	var preference models.Preference
	if result := initializers.DB.First(&preference, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	initializers.DB.Model(&preference).Updates(models.Preference{
		Name:   body.Name,
		Value:  body.Value,
		UserId: body.UserId,
	})

	c.JSON(http.StatusOK, gin.H{"data": &preference})
}

func DeletePreference(c *gin.Context) {
	id := c.Param("id")
	var preference models.Preference

	if result := initializers.DB.First(&preference, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
	}

	initializers.DB.Delete(&preference)

	c.JSON(http.StatusOK, gin.H{"data": "Preference deleted successfully"})
}
