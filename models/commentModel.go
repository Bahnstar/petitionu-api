package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Text       string
	Sentiment  string
	UserId     uint
	PetitionId uint
}
