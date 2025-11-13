package repository

import (
	"learn/fiber/pkg/model/entity"

	"gorm.io/gorm"
)

type BlogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

func (r *BlogRepository) Create(blog *entity.Blog) error {
	return r.db.Create(blog).Error
}
