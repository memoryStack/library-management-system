package models

import "gorm.io/gorm"

type Book struct {
	Title string `json:"title" gorm:"not null"`
	Author string `json:"author" gorm:"not null"`
	ISBN string `json:"isbn" gorm:"not null;unique"`
	PublicationYear int `json:"publication_year" gorm:"not null"`
	Publisher string `json:"publisher" gorm:"not null"`
	Genre string `json:"genre" gorm:"not null"`
	Price float64 `json:"price" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`
	ImageURL string `json:"image_url" gorm:"not null"`
	RelatedImages []string `json:"related_images" gorm:"type:jsonb;serializer:json"`
	gorm.Model
}
