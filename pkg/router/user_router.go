package router

import (
	"learn/fiber/pkg/handler"
	"learn/fiber/pkg/middleware"
	"learn/fiber/pkg/service"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(app fiber.Router) {
	userService := service.NewUserService()
	userHandler := handler.NewUserHandler(userService)

	user := app.Group("/user")
	user.Get("/:id", middleware.JWTMidleware, userHandler.FindByIdHandler)
}
