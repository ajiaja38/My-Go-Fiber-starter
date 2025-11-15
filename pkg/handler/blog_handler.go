package handler

import (
	"fmt"
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

// @Summary		    Find All Blogs Paginate
// @Description	Get a list of all Blogs with pagination
// @Tags			      Blog
// @Accept			     json
// @Produce		    json
// @Param			request	query	model.PaginationRequest	true		"Pagination Request Payload"
// @Success		 		 		200						{object}	model.ResponseEntityPagination[[]res.FindBlogResponse]
// @Failure		 		 		401						{object}	model.ResponseError[any]
// @Router			     /blog/paginate [get]
func (b *BlogHandler) FindAllPaginateHandler(c *fiber.Ctx) error {
	var params model.PaginationRequest

	if err := c.QueryParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if params.Page <= 0 {
		params.Page = 1
	}

	if params.Limit <= 0 {
		params.Limit = 5
	}

	meta, blogs, err := b.blogService.FindAllPaginate(&params)

	if err != nil {
		return err
	}

	return utils.SuccessResponsePaginate(
		c,
		fiber.StatusOK,
		"Success Find All Blogs Paginate",
		blogs,
		meta,
	)
}

// @Summary		    Find Blog By Id
// @Description	Get Blog details by ID
// @Tags			      Blog
// @Accept			     json
// @Produce		    json
// @Param			id	path	string	true		"blog ID"
// @Success		 	 		200		{object}	model.ResponseEntityPagination[res.FindBlogResponse]
// @Failure		 	 		401		{object}	model.ResponseError[any]
// @Failure		 	 		404		{object}	model.ResponseError[any]
// @Router			     /blog/{id} [get]
func (b *BlogHandler) FindBlogByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	blog, err := b.blogService.FindById(id)

	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fmt.Sprintf("Success Get blog %s", blog.Title), blog)
}
