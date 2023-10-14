package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email          string `gorm:"unique"`
	Password       string
	FirstName      string
	LastName       string
	OrganizationId uint
	Preferences    []Preference
	Petitions      []Petition `gorm:"many2many:user_petitions;"`
}
