package models

import "gorm.io/gorm"

type Petition struct {
	gorm.Model
	OwnerId        uint
	Name           string
	Description    string
	OrganizationId uint
	Comments       []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	// Users          []User    `gorm:"many2many:user_petitions;"`
}
