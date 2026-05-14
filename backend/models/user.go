package models

import "gorm.io/gorm"

type User struct {
	Name string `json:"name" gorm:"not null"`
	Email  string `json:"email" gorm:"not null;unique"`
	EmailVerified     bool `json:"email_verified"`
	Image     string `json:"image_url"`
	Auth0ID   string `json:"auth0_id" gorm:"not null;unique"`
	gorm.Model
}
