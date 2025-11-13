package handler

import (
	"learn/fiber/pkg/model"
	"learn/fiber/pkg/model/req"
	"learn/fiber/pkg/service"
	"learn/fiber/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type BlogHandler struct {
	blogService service.BlogService
	validator   *validator.Validate
}

func NewBlogHandler(blogService service.BlogService) *BlogHandler {
	return &BlogHandler{
		blogService: blogService,
		validator:   validator.New(),
	}
}

// @Summary		    Create Blog
// @Description	Create a new blog
// @Tags			       Blog
// @Accept			     json
// @Produce		    json
// @Security		        BearerAuth
// @Param			request	body	req.CreateBlogDto	true	"Create Blog Request Payload"
// @Router			     /blog [post]
func (b *BlogHandler) CreateBlogHandler(c *fiber.Ctx) error {
	var payload req.CreateBlogDto

	if err := utils.ValidateRequestBody(c, b.validator, &payload); err != nil {
		return err
	}

	blog, err := b.blogService.CreateBlog(&payload, c.Locals("payload").(model.JwtPayload).Id)

	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Succes Create Blog 🚀", blog)
}
