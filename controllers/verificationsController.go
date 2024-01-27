package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/likexian/whois"
	"github.com/likexian/whois-parser"

	"github.com/Bahnstar/petitionu-api/initializers"
	"github.com/Bahnstar/petitionu-api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

func GetOrganizationFromDNS(c *gin.Context) {
  email := "qc@students.llu.edu"
  email_domain := strings.Split(email, "@")
  split_email := strings.Split(email_domain[1], ".")
  domain := strings.Join(split_email[len(split_email) - 2:], ".")
  fmt.Println(domain)
	raw, err := whois.Whois(domain)
	if err != nil {
		panic(err)
	}

	parsed_whois, err := whoisparser.Parse(raw)
	if err != nil {
		panic(err)
	}
	fmt.Println("Organization: " + parsed_whois.Registrant.Organization)
}

func generateVerificationToken() (string, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

func SendVerification(email string) {
	var user models.User
	initializers.DB.First(&user, "email = ?", email)

	if user.ID == 0 {
		panic("Could not find user")
	}

	verification_token, err := generateVerificationToken()
	if err != nil {
		panic(err)
	}

	user.VerificationToken = verification_token
	user.EmailVerified = false
	if result := initializers.DB.Save(&user); result.Error != nil {
		panic("Could not update user with verification token")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "owners@bahnstar.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "PetitionU - Verify Email")
	m.SetBody("text/html", "Please click <a href='http://localhost:8080/verify?token="+verification_token+"'>here</a> to verify your email.")

	pw := os.Getenv("MAIL_PASSWORD")
	d := gomail.NewDialer("smtp.gmail.com", 587, "jacobabahn@gmail.com", pw)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

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

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
