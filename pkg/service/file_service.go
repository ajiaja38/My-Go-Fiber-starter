package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	env "learn/fiber/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

type FileService interface {
	Upload(file *multipart.FileHeader) (string, error)
}

type fileService struct {
	client *s3.Client
	bucket string
}

func NewFileService() (FileService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(env.S3_REGION.GetValue()),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				env.S3_ACCESS_KEY.GetValue(),
				env.S3_SECRET_KEY.GetValue(),
				"",
			),
		),
	)

	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(env.S3_ENDPOINT.GetValue())
		o.UsePathStyle = true
	})

	return &fileService{
		client: client,
		bucket: env.S3_BUCKET.GetValue(),
	}, nil
}

func (f *fileService) Upload(file *multipart.FileHeader) (string, error) {
	fileContent, err := file.Open()

	if err != nil {
		return "", fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	defer fileContent.Close()

	key := fmt.Sprintf("data/kehadiran/image/%s", time.Now().Format("20060102150405_")+file.Filename)

	uploader := manager.NewUploader(f.client)

	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:             aws.String(f.bucket),
		Key:                aws.String(key),
		Body:               fileContent,
		ACL:                "public-read",
		ContentType:        aws.String(file.Header.Get("Content-Type")),
		ContentDisposition: aws.String("inline"),
	})

	if err != nil {
		return "", fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	url := fmt.Sprintf("https://%s.nos.wjv-1.neo.id/%s", env.S3_BUCKET.GetValue(), key)

	return url, nil
}
