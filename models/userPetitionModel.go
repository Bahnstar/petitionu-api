package models

import "gorm.io/gorm"

type Relationship string

const (
	Owner  Relationship = "owner"
	Signee Relationship = "signee"
)

type UserPetition struct {
	gorm.Model
	UserID       int `gorm:"primaryKey"`
	PetitionID   int `gorm:"primaryKey"`
	Relationship Relationship
	Bookmarked   bool
}
