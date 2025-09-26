package handler

import (
	"learn/fiber/pkg/model"
	"learn/fiber/pkg/service"
	"learn/fiber/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator.New(),
	}
}

func (u *UserHandler) RegisterUserHandler(c *fiber.Ctx) error {
	var payload model.UserRegisterRequest

	if err := utils.ValidateRequestBody(c, u.validator, &payload); err != nil {
		return err
	}

	if !utils.ValidatePassword(payload.Password) {
		return fiber.NewError(fiber.StatusBadRequest, "Password must be at least 6 characters long, contain at least one uppercase letter, one number, and one special character")
	}

	user, err := u.userService.RegisterUser(payload)

	if err != nil {
		return err
	}

	return c.JSON(model.ResponseEntity[model.UserResponse]{
		Code:    fiber.StatusCreated,
		Message: "Succes Register User",
		Data:    *user,
	})
}

func (u *UserHandler) LoginUserHandler(c *fiber.Ctx) error {
	var payload model.UserLoginRequest

	if err := utils.ValidateRequestBody(c, u.validator, &payload); err != nil {
		return err
	}

	jwtResponse, err := u.userService.LoginUser(payload)

	if err != nil {
		return err
	}

	return c.JSON(model.ResponseEntity[model.JwtResponse]{
		Code:    fiber.StatusOK,
		Message: "Succes Login User",
		Data:    *jwtResponse,
	})
}

func (u *UserHandler) RefreshTokenHandler(c *fiber.Ctx) error {
	var payload model.RefreshTokenRequest

	if err := utils.ValidateRequestBody(c, u.validator, &payload); err != nil {
		return err
	}

	refreshTokenResponse, err := u.userService.RefreshToken(payload.RefreshToken)

	if err != nil {
		return err
	}

	return c.JSON(model.ResponseEntity[model.RefreshTokenResponse]{
		Code:    fiber.StatusOK,
		Message: "Succes Refresh Token",
		Data:    *refreshTokenResponse,
	})
}

func (u *UserHandler) FindAllHandler(c *fiber.Ctx) error {
	users, err := u.userService.FindAll()

	if err != nil {
		return err
	}

	return c.JSON(model.ResponseEntity[[]model.UserResponse]{
		Code:    fiber.StatusOK,
		Message: "Succes Find All User",
		Data:    users,
	})
}

func (u *UserHandler) FindByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := u.userService.FindById(id)

	if err != nil {
		return err
	}

	return c.JSON(model.ResponseEntity[model.UserResponse]{
		Code:    fiber.StatusOK,
		Message: "Succes Find User By Id",
		Data:    *user,
	})
}

func (u *UserHandler) UpdateUserByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var payload model.UserUpdateRequest

	if err := utils.ValidateRequestBody(c, u.validator, &payload); err != nil {
		return err
	}

	user, err := u.userService.UpdateUserById(id, payload)

	if err != nil {
		return err
	}

	return c.JSON(model.ResponseEntity[model.UserResponse]{
		Code:    fiber.StatusOK,
		Message: "Succes Update User By Id",
		Data:    *user,
	})
}

func (u *UserHandler) DeleteUserByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := u.userService.DeleteUserById(id); err != nil {
		return err
	}

	return c.JSON(model.ResponseEntity[any]{
		Code:    fiber.StatusOK,
		Message: "Succes Delete User By Id",
		Data:    nil,
	})
}
