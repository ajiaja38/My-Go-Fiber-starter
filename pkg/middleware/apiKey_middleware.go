package middleware

import (
	"learn/fiber/config"

	"github.com/gofiber/fiber/v2"
)

func ApiKeyGuard(c *fiber.Ctx) error {
	apiKey := c.Get("X-Api-Key")

	if apiKey != config.API_KEY.GetValue() {
		return fiber.NewError(fiber.StatusUnauthorized, "Error, Unauthorized!")
	}

	return c.Next()
}
