package router

import (
	"learn/fiber/pkg/handler"
	"learn/fiber/pkg/middleware"
	"learn/fiber/pkg/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRouter(app fiber.Router, db *gorm.DB) {
	userService := service.NewUserService(db)
	userHandler := handler.NewUserHandler(userService)

	user := app.Group("/user")

	user.Get("/", middleware.JWTMidleware, userHandler.FindAllHandler)
	user.Get("/:id", middleware.JWTMidleware, userHandler.FindByIdHandler)
}
