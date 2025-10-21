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

// @Summary		    Register User
// @Description	Register a new user
// @Tags			       user
// @Accept			     json
// @Produce		    json
// @Param			request	body	model.UserRegisterRequest	true		"Register Request Payload"
// @Success		 		 		201							{object}	model.ResponseEntity[model.UserResponse]
// @Router			     /user/register [post]
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

// @Summary		    Login User
// @Description	Log in a user
// @Tags			       user
// @Accept			     json
// @Produce		    json
// @Param			request	body	model.UserLoginRequest	true		"Login Request Payload"
// @Success		 		 		200						{object}	model.ResponseEntity[model.JwtResponse]
// @Router			     /user/login [post]
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

// @Summary		    Refresh Token
// @Description	Refresh JWT tokens
// @Tags			       user
// @Accept			     json
// @Produce		    json
// @Security		        BearerAuth
// @Param			request	body	model.RefreshTokenRequest	true		"Refresh Token Request Payload"
// @Success		 		 		200							{object}	model.ResponseEntity[model.RefreshTokenResponse]
// @Router			     /user/refresh-token [put]
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

// @Summary		    Find All Users
// @Description	Get a list of all users
// @Tags			       user
// @Accept			     json
// @Produce		    json
// @Security		        BearerAuth
// @Success		 	 	200	{object}	model.ResponseEntity[[]model.UserResponse]
// @Failure		 	 	401	{object}	model.ResponseError[any]
// @Router			     /user [get]
func (u *UserHandler) FindAllHandler(c *fiber.Ctx) error {
	users, err := u.userService.FindAll()

	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Succes Find All Users", users)
}

// @Summary		    Find User By Id
// @Description	Get user details by ID
// @Tags			       user
// @Accept			     json
// @Produce		    json
// @Security		        BearerAuth
// @Param			id	path	string	true		"User ID"
// @Success		 	 		200		{object}	model.ResponseEntity[model.UserResponse]
// @Failure		 	 		401		{object}	model.ResponseError[any]
// @Failure		 	 		404		{object}	model.ResponseError[any]
// @Router			     /user/{id} [get]
func (u *UserHandler) FindByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := u.userService.FindById(id)

	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Succes Find User By Id", user)
}

// @Summary		    Update User By Id
// @Description	Update user details by ID
// @Tags			       user
// @Accept			     json
// @Produce		    json
// @Security		        BearerAuth
// @Param			request	body	model.UserUpdateRequest	true		"Update User Request Payload"
// @Param			id		path	string					true		"User ID"
// @Success		 		 		200						{object}	model.ResponseEntity[model.UserResponse]
// @Router			     /user/{id} [put]
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

// @Summary		    Delete User By Id
// @Description	Delete user details by ID
// @Tags			       user
// @Accept			     json
// @Produce		    json
// @Security		        BearerAuth
// @Param			id	path	string	true		"User ID"
// @Success		 	 		200		{object}	model.ResponseEntity[any]
// @Failure		 	 		401		{object}	model.ResponseError[any]
// @Failure		 	 		404		{object}	model.ResponseError[any]
// @Router			    /user/{id} [delete]
func (u *UserHandler) DeleteUserByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := u.userService.DeleteUserById(id); err != nil {
		return err
	}

	return utils.SuccessResponse[*struct{}](c, fiber.StatusOK, "Succes Delete User By Id", nil)
}
