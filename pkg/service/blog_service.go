package service

import (
	"learn/fiber/pkg/model"
	"learn/fiber/pkg/model/entity"
	"learn/fiber/pkg/model/req"
	"learn/fiber/pkg/model/res"
	"learn/fiber/pkg/repository"

	"github.com/gofiber/fiber/v2"
)

type BlogService interface {
	CreateBlog(createBlogDto *req.CreateBlogDto, userId string) (*entity.Blog, error)
	FindAllPaginate(pagination *model.PaginationRequest) (*model.MetaPagination, *[]res.FindBlogResponse, error)
	FindById(id string) (*res.FindBlogResponse, error)
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

func (b *blogService) FindAllPaginate(pagination *model.PaginationRequest) (*model.MetaPagination, *[]res.FindBlogResponse, error) {
	blogs, total, err := b.repository.FindAllPagination(pagination.Page, pagination.Limit, pagination.Search)

	if err != nil {
		return nil, nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	totalPage := (total + int64(pagination.Limit) - 1) / int64(pagination.Limit)

	meta := &model.MetaPagination{
		Page:      pagination.Page,
		Limit:     pagination.Limit,
		TotalPage: int(totalPage),
		TotalData: int(total),
	}

	return meta, blogs, nil
}

func (b *blogService) FindById(id string) (*res.FindBlogResponse, error) {
	blog, err := b.repository.FindById(id)

	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return blog, nil
}
