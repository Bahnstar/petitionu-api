package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/Bahnstar/petitionu-api/initializers"
	"github.com/Bahnstar/petitionu-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserBody struct {
	Email          string `binding:"required"`
	Password       string `binding:"required"`
	FirstName      string
	LastName       string
	OrganizationId uint `binding:"required"`
	Preferences    []models.Preference
	Petitions      []models.Petition
}

type UpdateUserBody struct {
	Email          string
	Password       string
	FirstName      string
	LastName       string
	OrganizationId uint
	Preferences    []models.Preference
	Petitions      []models.Petition
}

func SignUp(c *gin.Context) {
	var body CreateUserBody

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	var organization models.Organization

	if organization_result := initializers.DB.First(&organization, body.OrganizationId); organization_result.Error != nil {
		c.AbortWithError(http.StatusNotFound, organization_result.Error)
		return
	}

	user := models.User{Email: body.Email, Password: string(hash), OrganizationId: organization.ID}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	SendVerification(user.Email)

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})

}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{"message": user})
}

func GetUsers(c *gin.Context) {
	var users []models.User

	if result := initializers.DB.Preload("Preferences").Preload("Petitions").Find(&users); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if result := initializers.DB.Preload("Preferences").Preload("Petitions").First(&user, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var body UpdateUserBody

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body in update user"})
	}

	var user models.User

	if result := initializers.DB.First(&user, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	initializers.DB.Model(&user).Updates(models.User{
		Email:          body.Email,
		Password:       body.Password,
		FirstName:      body.FirstName,
		LastName:       body.LastName,
		OrganizationId: body.OrganizationId,
		Preferences:    body.Preferences,
		Petitions:      body.Petitions,
	})

	c.JSON(http.StatusOK, gin.H{"data": &user})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if result := initializers.DB.First(&user, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
	}

	initializers.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"data": "User deleted successfully"})
}
