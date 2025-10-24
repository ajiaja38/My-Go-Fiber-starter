package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func ApiKeyGuard(c *fiber.Ctx) error {
	apiKey := c.Get("X-Api-Key")

	if apiKey != os.Getenv("API_KEY") {
		return fiber.NewError(fiber.StatusUnauthorized, "Error, Unauthorized!")
	}

	return c.Next()
}
