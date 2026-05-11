package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"library-management-system/backend/initializers"
	"library-management-system/backend/models"
)

func CreateBook(c *fiber.Ctx) error {
	fmt.Println("Creating book")
	// book := c.Bind()
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result := initializers.DB.Create(&book) // pass pointer of data to Create

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Book created successfully",
		"book": book,
	})
}
