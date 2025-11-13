package service

import (
	"learn/fiber/pkg/model/entity"
	"learn/fiber/pkg/model/req"
	"learn/fiber/pkg/repository"

	"github.com/gofiber/fiber/v2"
)

type BlogService interface {
	CreateBlog(createBlogDto *req.CreateBlogDto, userId string) (*entity.Blog, error)
}

type blogService struct {
	repository     *repository.BlogRepository
	userRepository *repository.UserRepository
}

func NewBlogService(repository *repository.BlogRepository, userRepository *repository.UserRepository) BlogService {
	return &blogService{repository: repository, userRepository: userRepository}
}

func (b *blogService) CreateBlog(createBlogDto *req.CreateBlogDto, userId string) (*entity.Blog, error) {
	user, err := b.userRepository.FindById(userId)

	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	blog := &entity.Blog{
		Title:  createBlogDto.Title,
		Body:   createBlogDto.Body,
		Image:  createBlogDto.Image,
		UserId: user.Id,
	}

	if err := b.repository.Create(blog); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return blog, nil
}
