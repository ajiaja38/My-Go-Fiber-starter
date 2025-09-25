package router

import "github.com/gofiber/fiber/v2"

func MainRouter(app fiber.Router) {
	UserRouter(app)
}
