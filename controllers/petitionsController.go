package controllers

import (
	"net/http"

	"github.com/Bahnstar/petitionu-api/initializers"
	"github.com/Bahnstar/petitionu-api/models"
	"github.com/gin-gonic/gin"
)

type CreatePetitionBody struct {
	OwnerId        uint   `binding:"required"`
	Name           string `binding:"required"`
	Description    string
	OrganizationId uint `binding:"required"`
	Comments       []models.Comment
}

type UpdatePetitionBody struct {
	OwnerId        uint
	Name           string
	Description    string
	OrganizationId uint
	Comments       []models.Comment
}

func GetPetitions(c *gin.Context) {
	var petitions []models.Petition

	if result := initializers.DB.Statement.Preload("Comments").Find(&petitions); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
	}

	c.JSON(http.StatusOK, gin.H{"data": &petitions})
}

func GetPetition(c *gin.Context) {
	id := c.Param("id")
	var petition models.Petition

	if result := initializers.DB.Statement.Preload("Comments").First(&petition, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": &petition})
}

func CreatePetition(c *gin.Context) {
	var body CreatePetitionBody

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	var user models.User

	if result := initializers.DB.First(&user, int(body.OwnerId)); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
	}

	var organization models.Organization

	if result := initializers.DB.First(&organization, body.OrganizationId); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	petition := models.Petition{
		OwnerId:        body.OwnerId,
		Name:           body.Name,
		Description:    body.Description,
		OrganizationId: organization.ID,
		Comments:       body.Comments,
	}

	result := initializers.DB.Create(&petition)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create petition"})
		return
	}

	userPetition := models.UserPetition{
		UserID:       int(user.ID),
		PetitionID:   int(petition.ID),
		Relationship: models.Owner,
		Bookmarked:   false,
	}

	result = initializers.DB.Create(&userPetition)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add petition in user_petitions join table"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Petition created successfully"})
}

func UpdatePetition(c *gin.Context) {
	id := c.Param("id")
	var body UpdatePetitionBody

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body in update petition"})
	}

	var petition models.Petition
	if result := initializers.DB.First(&petition, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	var organization models.Organization
	if body.OrganizationId != 0 {
		if result := initializers.DB.First(&organization, body.OrganizationId); result.Error != nil {
			c.AbortWithError(http.StatusNotFound, result.Error)
			return
		}
	}

	initializers.DB.Model(&petition).Updates(models.Petition{
		OwnerId:        body.OwnerId,
		Name:           body.Name,
		Description:    body.Description,
		OrganizationId: body.OrganizationId,
		Comments:       body.Comments,
	})

	c.JSON(http.StatusOK, gin.H{"data": &petition})
}

func DeletePetition(c *gin.Context) {
	id := c.Param("id")
	var petition models.Petition

	if result := initializers.DB.First(&petition, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
	}

	initializers.DB.Delete(&petition)

	c.JSON(http.StatusOK, gin.H{"data": "Petition deleted successfully"})
}
