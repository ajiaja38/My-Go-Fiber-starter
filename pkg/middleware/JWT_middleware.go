package middleware

import (
	"learn/fiber/utils"
	"os"

	"github.com/gofiber/fiber/v2"
)

func JWTMidleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized, no token provided")
	}

	tokenStr := authHeader[len("Bearer "):]

	payload, err := utils.ValidateToken(tokenStr, os.Getenv("JWT_SECRET_ACCESS_TOKEN"))

	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	c.Locals("payload", payload)

	return c.Next()
}
