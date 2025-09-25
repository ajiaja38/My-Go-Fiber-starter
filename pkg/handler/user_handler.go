package handler

import (
	"learn/fiber/pkg/model"
	"learn/fiber/pkg/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) FindByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := u.userService.FindById(id)

	if err != nil {
		return err
	}

	return c.JSON(model.ResponseEntity[model.User]{
		Code:    fiber.StatusOK,
		Message: "success",
		Data:    *user,
	})
}
