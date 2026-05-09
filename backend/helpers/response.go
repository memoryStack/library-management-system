package helpers

import "github.com/gofiber/fiber/v2"

// JSONError sends a consistent error payload.
func JSONError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}
