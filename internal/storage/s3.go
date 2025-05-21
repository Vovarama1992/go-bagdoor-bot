package storage

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Uploader struct {
	Client   *s3.Client
	Bucket   string
	Endpoint string
}

func NewS3Uploader() *S3Uploader {
	region := os.Getenv("S3_REGION")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	bucket := os.Getenv("S3_BUCKET_NAME")
	endpoint := os.Getenv("S3_ENDPOINT")

	if region == "" || accessKey == "" || secretKey == "" || bucket == "" || endpoint == "" {
		log.Fatal("S3 env-переменные не заданы полностью")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		log.Fatalf("не удалось инициализировать AWS config: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3Uploader{
		Client:   client,
		Bucket:   bucket,
		Endpoint: endpoint,
	}
}

func (u *S3Uploader) Upload(orderID int, fileName string, content []byte) (string, error) {
	key := fmt.Sprintf("orders/%d/%s", orderID, fileName)
	contentType := mime.TypeByExtension(filepath.Ext(fileName))

	_, err := u.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(u.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(contentType),
		ACL:         s3types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s/%s", u.Endpoint, u.Bucket, key)
	return url, nil
}
