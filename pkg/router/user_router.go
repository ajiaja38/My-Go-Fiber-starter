package router

import (
	"learn/fiber/pkg/enum"
	"learn/fiber/pkg/handler"
	"learn/fiber/pkg/middleware"
	"learn/fiber/pkg/repository"
	"learn/fiber/pkg/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRouter(app fiber.Router, db *gorm.DB) {
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	user := app.Group("/user")

	user.Post("/register", userHandler.RegisterUserHandler)
	user.Post("/login", userHandler.LoginUserHandler)
	user.Get("/", middleware.JWTMidleware, middleware.RoleMiddleware(enum.ROLE_USER, enum.ROLE_ADMIN), userHandler.FindAllHandler)
	user.Get("/paginate", middleware.JWTMidleware, middleware.RoleMiddleware(enum.ROLE_USER, enum.ROLE_ADMIN), userHandler.FindAllPaginateHandler)
	user.Get("/:id", middleware.JWTMidleware, userHandler.FindByIdHandler)
	user.Put("/refresh-token", userHandler.RefreshTokenHandler)
	user.Put("/:id", middleware.JWTMidleware, userHandler.UpdateUserByIdHandler)
	user.Delete("/:id", middleware.JWTMidleware, middleware.RoleMiddleware(enum.ROLE_ADMIN), userHandler.DeleteUserByIdHandler)
}
