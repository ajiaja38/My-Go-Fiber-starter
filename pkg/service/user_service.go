package service

import (
	"learn/fiber/pkg/model"
	"learn/fiber/pkg/repository"
	"learn/fiber/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(payload model.UserRegisterRequest) (*model.UserResponse, error)
	LoginUser(payload model.UserLoginRequest) (*model.JwtResponse, error)
	RefreshToken(refreshToken string) (*model.RefreshTokenResponse, error)
	FindAll() ([]model.UserResponse, error)
	FindAllPaginated(pagination model.PaginationRequest) (*model.MetaPagination, []model.UserResponse, error)
	FindById(id string) (*model.UserResponse, error)
	UpdateUserById(id string, payload model.UserUpdateRequest) (*model.UserResponse, error)
	DeleteUserById(id string) error
}

type userService struct {
	repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) UserService {
	return &userService{
		repository: repository,
	}
}

func (u *userService) RegisterUser(payload model.UserRegisterRequest) (*model.UserResponse, error) {
	if payload.Password != payload.ConfirmPassword {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Password and Confirm Password do not match")
	}

	passwordHashed, err := hashedPassword(payload.Password)

	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user := model.User{
		Email:    payload.Email,
		Username: payload.Username,
		Password: passwordHashed,
		Role:     payload.Role,
	}

	if err := u.repository.Create(&user); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userResponse := transformUserResponse(user)

	return &userResponse, nil
}

func (u *userService) LoginUser(payload model.UserLoginRequest) (*model.JwtResponse, error) {
	user, err := u.repository.FindByEmail(payload.Email)

	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	if !checkPasswordHash(payload.Password, user.Password) {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	jwtPayload := model.JwtPayload{
		Id:   user.Id,
		Role: user.Role,
	}

	accessToken, err := utils.GenerateAccessToken(jwtPayload)

	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	refreshToken, err := utils.GenerateRefreshToken(jwtPayload)

	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return &model.JwtResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *userService) RefreshToken(refreshToken string) (*model.RefreshTokenResponse, error) {
	accessToken, err := utils.GenerateNewAccessToken(refreshToken)

	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return &model.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (u *userService) FindAll() ([]model.UserResponse, error) {
	users, err := u.repository.FindAll()

	if err != nil {
		return nil, err
	}

	var userResponses []model.UserResponse

	for _, user := range users {
		userResponse := transformUserResponse(user)
		userResponses = append(userResponses, userResponse)
	}

	return userResponses, nil
}

func (u *userService) FindAllPaginated(pagination model.PaginationRequest) (*model.MetaPagination, []model.UserResponse, error) {
	users, totalData, err := u.repository.FindAllPaginated(pagination.Page, pagination.Limit, pagination.Search)

	if err != nil {
		return nil, nil, err
	}

	var userResponses []model.UserResponse = []model.UserResponse{}

	for _, user := range users {
		userResponse := transformUserResponse(user)
		userResponses = append(userResponses, userResponse)
	}

	totalPage := (totalData + int64(pagination.Limit) - 1) / int64(pagination.Limit)

	meta := &model.MetaPagination{
		Page:      pagination.Page,
		Limit:     pagination.Limit,
		TotalPage: int(totalPage),
		TotalData: int(totalData),
	}

	return meta, userResponses, nil
}

func (u *userService) FindById(id string) (*model.UserResponse, error) {
	user, err := u.repository.FindById(id)

	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	userResponse := transformUserResponse(*user)

	return &userResponse, nil
}

func (u *userService) UpdateUserById(id string, payload model.UserUpdateRequest) (*model.UserResponse, error) {
	user, err := u.repository.FindById(id)

	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	user.Email = payload.Email
	user.Username = payload.Username
	user.Role = payload.Role

	if err := u.repository.Update(user); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userResponse := transformUserResponse(*user)

	return &userResponse, nil
}

func (u *userService) DeleteUserById(id string) error {
	if err := u.repository.Delete(id); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return nil
}

func transformUserResponse(user model.User) model.UserResponse {
	userResponse := model.UserResponse{
		Id:        user.Id,
		Email:     user.Email,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.DeletedAt.Valid {
		deleted := user.DeletedAt.Time
		userResponse.DeletedAt = &deleted
	}

	return userResponse
}

func hashedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
