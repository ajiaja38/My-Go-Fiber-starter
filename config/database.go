package config

import (
	"fmt"
	"learn/fiber/pkg/model/entity"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBConfig() *gorm.DB {
	host := DB_HOST.GetValue()
	user := DB_USER.GetValue()
	password := DB_PASSWORD.GetValue()
	dbname := DB_NAME.GetValue()
	port := DB_PORT.GetValue()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	AutoMigrateEntity(db)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

func AutoMigrateEntity(db *gorm.DB) {
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Blog{})
}
