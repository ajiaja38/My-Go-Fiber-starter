package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MainRouter(app fiber.Router, db *gorm.DB) {
	UserRouter(app, db)
}
