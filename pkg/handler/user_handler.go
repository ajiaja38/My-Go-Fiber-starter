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

	return utils.SuccessResponse(c, fiber.StatusCreated, "Succes Register User", user)
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

	return utils.SuccessResponse(c, fiber.StatusOK, "Succes Login User", jwtResponse)
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

	return utils.SuccessResponse(c, fiber.StatusOK, "Succes Refresh Token", refreshTokenResponse)
}

func (u *UserHandler) FindAllHandler(c *fiber.Ctx) error {
	users, err := u.userService.FindAll()

	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Succes Find All Users", users)
}

func (u *UserHandler) FindByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := u.userService.FindById(id)

	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Succes Find User By Id", user)
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

	return utils.SuccessResponse(c, fiber.StatusOK, "Succes Update User By Id", user)
}

func (u *UserHandler) DeleteUserByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := u.userService.DeleteUserById(id); err != nil {
		return err
	}

	return utils.SuccessResponse[*struct{}](c, fiber.StatusOK, "Succes Delete User By Id", nil)
}
