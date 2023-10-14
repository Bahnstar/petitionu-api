package models

import "gorm.io/gorm"

type Preference struct {
	gorm.Model
	Name   string
	Value  string
	UserId uint
}
