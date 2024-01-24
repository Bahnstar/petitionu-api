package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email             string `gorm:"unique"`
	Password          string
	FirstName         string
	LastName          string
	OrganizationId    uint
	EmailVerified     bool
	VerificationToken string
	Preferences       []Preference `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Petitions         []Petition   `gorm:"many2many:user_petitions;"`
}
