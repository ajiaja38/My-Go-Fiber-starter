package service

import (
	"fmt"
	"learn/fiber/pkg/model"
	"learn/fiber/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	RegisterUser(payload model.UserRegisterRequest) (*model.UserResponse, error)
	LoginUser(payload model.UserLoginRequest) (*model.JwtResponse, error)
	RefreshToken(refreshToken string) (*model.RefreshTokenResponse, error)
	FindAll() ([]model.UserResponse, error)
	FindById(id string) (*model.UserResponse, error)
	UpdateUserById(id string, payload model.UserUpdateRequest) (*model.UserResponse, error)
	DeleteUserById(id string) error
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
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

	if err := u.db.Create(&user).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userResponse := transformUserResponse(user)

	return &userResponse, nil
}

func (u *userService) LoginUser(payload model.UserLoginRequest) (*model.JwtResponse, error) {
	var user model.User

	if err := u.db.First(&user, "email = ?", payload.Email).Error; err != nil {
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
	var users []model.User

	if err := u.db.Find(&users).Error; err != nil {
		return nil, err
	}

	var userResponses []model.UserResponse

	for _, user := range users {
		userResponse := transformUserResponse(user)
		userResponses = append(userResponses, userResponse)
	}

	return userResponses, nil
}

func (u *userService) FindById(id string) (*model.UserResponse, error) {
	var user model.User

	if err := u.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("User not found: %s", err.Error()))
	}

	userResponse := transformUserResponse(user)

	return &userResponse, nil
}

func (u *userService) UpdateUserById(id string, payload model.UserUpdateRequest) (*model.UserResponse, error) {
	var user model.User

	if err := u.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("User not found: %s", err.Error()))
	}

	user.Email = payload.Email
	user.Username = payload.Username
	user.Role = payload.Role

	if err := u.db.Save(&user).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userResponse := transformUserResponse(user)

	return &userResponse, nil
}

func (u *userService) DeleteUserById(id string) error {
	if err := u.db.Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("User not found: %s", err.Error()))
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
