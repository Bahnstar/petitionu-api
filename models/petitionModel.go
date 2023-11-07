package models

import "gorm.io/gorm"

type Petition struct {
	gorm.Model
	OwnerId        int
	Name           string
	Description    string
	OrganizationId uint
	Comments       []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
