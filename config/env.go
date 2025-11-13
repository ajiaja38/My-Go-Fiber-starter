package config

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvKey string

const (
	// General
	PORT    EnvKey = "PORT"
	API_KEY EnvKey = "API_KEY"

	// JWT
	JWT_SECRET_ACCESS_TOKEN  EnvKey = "JWT_SECRET_ACCESS_TOKEN"
	JWT_SECRET_REFRESH_TOKEN EnvKey = "JWT_SECRET_REFRESH_TOKEN"

	// Database
	DB_HOST     EnvKey = "DB_HOST"
	DB_USER     EnvKey = "DB_USER"
	DB_PASSWORD EnvKey = "DB_PASSWORD"
	DB_NAME     EnvKey = "DB_NAME"
	DB_PORT     EnvKey = "DB_PORT"

	// S3
	S3_ACCESS_KEY EnvKey = "S3_ACCESS_KEY"
	S3_SECRET_KEY EnvKey = "S3_SECRET_KEY"
	S3_BUCKET     EnvKey = "S3_BUCKET"
	S3_ENDPOINT   EnvKey = "S3_ENDPOINT"
	S3_REGION     EnvKey = "S3_REGION"
	S3_SERVE_URL  EnvKey = "S3_SERVE_URL"
)

func LoadEnv() error {
	return godotenv.Load()
}

func (e EnvKey) GetValue() string {
	return os.Getenv(string(e))
}
