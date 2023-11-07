package controllers

import (
	"net/http"

	"github.com/Bahnstar/petitionu-api/initializers"
	"github.com/Bahnstar/petitionu-api/models"
	"github.com/gin-gonic/gin"
)

type CreateCommentBody struct {
	Text       string `binding:"required"`
	Sentiment  string
	UserId     uint `binding:"required"`
	PetitionId uint `binding:"required"`
}

type UpdateCommentBody struct {
	Text      string
	Sentiment string
}

func GetComments(c *gin.Context) {
	var comments []models.Comment

	if result := initializers.DB.Find(&comments); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
	}

	c.JSON(http.StatusOK, gin.H{"data": &comments})
}

func GetComment(c *gin.Context) {
	id := c.Param("id")
	var comment models.Comment

	if result := initializers.DB.First(&comment, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": &comment})
}

func CreateComment(c *gin.Context) {
	var body CreateCommentBody

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	var user models.User

	if result := initializers.DB.First(&user, body.UserId); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	var petition models.Petition

	if result := initializers.DB.First(&petition, body.PetitionId); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	comment := models.Comment{
		Text:       body.Text,
		Sentiment:  body.Sentiment,
		UserId:     body.UserId,
		PetitionId: body.PetitionId,
	}

	result := initializers.DB.Create(&comment)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully"})
}

func UpdateComment(c *gin.Context) {
	id := c.Param("id")
	var body UpdateCommentBody

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body in update comment"})
	}

	var comment models.Comment

	if result := initializers.DB.First(&comment, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	initializers.DB.Model(&comment).Updates(models.Comment{
		Text:      body.Text,
		Sentiment: body.Sentiment,
	})

	c.JSON(http.StatusOK, gin.H{"data": &comment})
}

func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	var comment models.Comment

	if result := initializers.DB.First(&comment, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
	}

	initializers.DB.Delete(&comment)

	c.JSON(http.StatusOK, gin.H{"data": "Comment deleted successfully"})
}
