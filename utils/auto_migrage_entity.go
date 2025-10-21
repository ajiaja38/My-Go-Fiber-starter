package utils

import (
	"learn/fiber/pkg/model"

	"gorm.io/gorm"
)

func AutoMigrateEntity(db *gorm.DB) {
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Blog{})
}
