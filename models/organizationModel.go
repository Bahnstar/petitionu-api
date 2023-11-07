package models

import "gorm.io/gorm"

type Organization struct {
	gorm.Model
	Name        string
	Description string
	Users       []User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Petitions   []Petition `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
