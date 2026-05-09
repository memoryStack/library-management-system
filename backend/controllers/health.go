package controllers

import "github.com/gofiber/fiber/v2"

// Health responds for load balancers and uptime checks.
func Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
