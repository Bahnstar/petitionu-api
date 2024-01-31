package helpers

import (
	"strings"
	"os"

	"github.com/Bahnstar/petitionu-api/models"
	"github.com/Bahnstar/petitionu-api/initializers"
	"github.com/likexian/whois"
	"github.com/likexian/whois-parser"
	"gopkg.in/gomail.v2"
	"github.com/google/uuid"
)

func GetDomainFromEmail(email string) string {
  email_domain := strings.Split(email, "@")
	split_email := strings.Split(email_domain[1], ".")
	domain := strings.Join(split_email[len(split_email)-2:], ".")

  return domain
}

func GetOrganizationFromDNS(domain string) string {
	raw, err := whois.Whois(domain)
	if err != nil {
		panic(err)
	}

	parsed_whois, err := whoisparser.Parse(raw)
	if err != nil {
		panic(err)
	}
	organiztion := parsed_whois.Registrant.Organization
	return organiztion
}

func SendVerificationEmail(email string) {
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

func generateVerificationToken() (string, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return token.String(), nil
}
