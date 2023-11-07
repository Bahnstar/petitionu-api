package models

import "gorm.io/gorm"

type Petition struct {
	gorm.Model
	ownerId        int
	Name           string
	Description    string
	OrganizationId uint
	Comments       []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
