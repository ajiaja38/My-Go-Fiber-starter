package router

import (
	"learn/fiber/pkg/handler"

	"github.com/gofiber/fiber/v2"
)

func FileRouter(app fiber.Router, fileHandler *handler.FileHandler) {
	file := app.Group("/file")

	file.Post("/upload", fileHandler.UploadFileHandler)
}
