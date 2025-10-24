package middleware

import "github.com/gofiber/fiber/v2"

func LimitUploadSize() fiber.Handler {
	return func(c *fiber.Ctx) error {

		if c.Request().Header.ContentLength() > 5*1024*1024 {
			return fiber.NewError(fiber.StatusBadRequest, "File size is too large")
		}

		return c.Next()
	}
}
