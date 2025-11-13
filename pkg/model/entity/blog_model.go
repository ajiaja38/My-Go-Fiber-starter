package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model
	Id     string `gorm:"primary_key" json:"id"`
	Title  string `gorm:"type:varchar(255); not null;" json:"title"`
	Body   string `gorm:"type:text; not null;" json:"body"`
	UserId string `gorm:"type:varchar(255); not null;" json:"userId"`
	User   User   `gorm:"foreignKey:UserId" json:"user"`
}

func (blog *Blog) BeforeCreate(db *gorm.DB) error {
	blog.Id = "blog-" + uuid.New().String()
	return nil
}
