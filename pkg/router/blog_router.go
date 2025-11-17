package router

import (
	"learn/fiber/pkg/handler"
	"learn/fiber/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func BlogRouter(app fiber.Router, blogHandler *handler.BlogHandler) {

	blog := app.Group("/blog")

	blog.Post("/", middleware.JWTMidleware, blogHandler.CreateBlogHandler)
	blog.Get("/paginate", blogHandler.FindAllPaginateHandler)
	blog.Get("/:id", blogHandler.FindBlogByIdHandler)

}
