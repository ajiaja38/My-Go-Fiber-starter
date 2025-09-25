package service

import (
	"learn/fiber/pkg/model"

	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	FindById(id string) (*model.User, error)
}

type userService struct{}

func NewUserService() UserService {
	return &userService{}
}

func (u *userService) FindById(id string) (*model.User, error) {
	user := model.User{
		Id:       id,
		Username: "Aji",
		Email:    "Nw9oQ@example.com",
	}

	err := false

	if err {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return &user, nil
}
