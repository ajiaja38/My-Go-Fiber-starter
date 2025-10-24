package router

import (
	"learn/fiber/pkg/handler"
	"learn/fiber/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func FileRouter(app fiber.Router, fileHandler *handler.FileHandler) {
	file := app.Group("/file")

	file.Post("/upload", middleware.ApiKeyGuard, fileHandler.UploadFileHandler)
}
