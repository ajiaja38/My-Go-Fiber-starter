package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	env "learn/fiber/config"
	"learn/fiber/pkg/model/res"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

type FileService interface {
	Upload(file *multipart.FileHeader) (*res.UploadFileResponse, error)
	Serve(s3Key string) (*s3.GetObjectOutput, error)
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

func (f *fileService) Upload(file *multipart.FileHeader) (*res.UploadFileResponse, error) {
	fileContent, err := file.Open()

	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	defer fileContent.Close()

	key := time.Now().Format("20060102150405") + filepath.Ext(file.Filename)

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
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	url := fmt.Sprintf("%s/%s", env.S3_SERVE_URL.GetValue(), key)

	urlResponse := res.UploadFileResponse{
		Url:      url,
		FileName: key,
	}

	return &urlResponse, nil
}

func (f *fileService) Serve(s3Key string) (*s3.GetObjectOutput, error) {
	getObjectRequest := &s3.GetObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(s3Key),
	}

	resp, err := f.client.GetObject(context.TODO(), getObjectRequest)

	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return resp, nil
}
