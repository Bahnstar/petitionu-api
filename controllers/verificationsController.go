package controllers

import (
	"fmt"
	"net/http"
	"github.com/Bahnstar/petitionu-api/initializers"
	"github.com/Bahnstar/petitionu-api/models"
	"github.com/gin-gonic/gin"
)


func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	fmt.Println("Token: ", token)

	var user models.User
	if result := initializers.DB.Where("verification_token = ?", token).First(&user); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	user.EmailVerified = true
	initializers.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": user.Email + "has been verified successfully"})
}
