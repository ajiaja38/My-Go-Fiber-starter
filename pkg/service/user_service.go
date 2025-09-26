package service

import (
	"learn/fiber/pkg/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserService interface {
	FindAll() ([]model.UserResponse, error)
	FindById(id string) (*model.UserResponse, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
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
		return nil, fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	userResponse := transformUserResponse(user)

	return &userResponse, nil
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
