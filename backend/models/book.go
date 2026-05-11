package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title string `gorm:"not null"`
	Author string `gorm:"not null"`
	ISBN string `gorm:"not null;unique"`
	PublicationYear int `gorm:"not null"`
	Publisher string `gorm:"not null"`
	Genre string `gorm:"not null"`
	Price float64 `gorm:"not null"`
	Description string `gorm:"not null"`
	ImageURL string `gorm:"not null"`
	RelatedImages []string `gorm:"type:jsonb"`
}
