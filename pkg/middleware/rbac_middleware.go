package middleware

import (
	"learn/fiber/pkg/enum"
	"learn/fiber/pkg/model"
	"slices"

	"github.com/gofiber/fiber/v2"
)

func RoleMiddleware(requiredRole ...enum.ERole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("payload").(model.JwtPayload).Role

		if slices.Contains(requiredRole, role) {
			return c.Next()
		}

		return fiber.NewError(fiber.StatusForbidden, "Forbidden Access, you don't have permission")
	}
}
