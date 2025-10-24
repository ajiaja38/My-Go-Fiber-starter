package service

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
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
		config.WithRegion(os.Getenv("S3_REGION")),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				os.Getenv("S3_ACCESS_KEY"),
				os.Getenv("S3_SECRET_KEY"),
				"",
			),
		),
	)

	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("S3_ENDPOINT"))
		o.UsePathStyle = true
	})

	return &fileService{
		client: client,
		bucket: os.Getenv("S3_BUCKET"),
	}, nil
}

func (f *fileService) Upload(file *multipart.FileHeader) (string, error) {
	fileContent, err := file.Open()

	if err != nil {
		return "", fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	defer fileContent.Close()

	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(fileContent)

	if err != nil {
		return "", fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	key := fmt.Sprintf("data/kehadiran/image/%s", time.Now().Format("20060102150405_")+file.Filename)

	_, err = f.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:             aws.String(f.bucket),
		Key:                aws.String(key),
		Body:               bytes.NewReader(buffer.Bytes()),
		ACL:                "public-read",
		ContentType:        aws.String(file.Header.Get("Content-Type")),
		ContentDisposition: aws.String("inline"),
	})

	if err != nil {
		return "", fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	url := fmt.Sprintf("https://%s.nos.wjv-1.neo.id/%s", os.Getenv("S3_BUCKET"), key)

	return url, nil
}
