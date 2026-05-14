package models

import "gorm.io/gorm"

type Author struct {
	FirstName string `json:"title" gorm:"not null"`
	LastName  string `json:"author" gorm:"not null"`
	Email     string `json:"email" gorm:"not null;unique"`
	About     string `json:"description" gorm:"not null"`
	Image     string `json:"image_url" gorm:"not null"`
	gorm.Model
}
