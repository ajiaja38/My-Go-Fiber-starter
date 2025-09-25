package err

import (
	"errors"
	"learn/fiber/pkg/model"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	var e *fiber.Error

	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	response := model.ResponseError[any]{
		ResponseEntity: model.ResponseEntity[any]{
			Code:    code,
			Message: message,
		},
		Path: c.Path(),
	}

	return c.Status(code).JSON(response)
}
