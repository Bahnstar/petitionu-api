package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Bahnstar/petitionu-api/initializers"
	"github.com/Bahnstar/petitionu-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		}

		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		}

		c.Set("user", user)

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
