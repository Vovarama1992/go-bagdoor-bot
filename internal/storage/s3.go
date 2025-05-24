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
	Client    *s3.Client
	Bucket    string
	PublicURL string
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
		Client:    client,
		Bucket:    bucket,
		PublicURL: fmt.Sprintf("%s/%s", endpoint, bucket),
	}
}

func (u *S3Uploader) UploadOrderMedia(orderID int, fileName string, content []byte) (string, error) {
	key := fmt.Sprintf("bot/orders/%d/%s", orderID, fileName)
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

	url := fmt.Sprintf("%s/%s", u.PublicURL, key)
	return url, nil
}

func (u *S3Uploader) UploadFlightMap(flightID int, fileName string, content []byte) (string, error) {
	key := fmt.Sprintf("bot/flights/%d/%s", flightID, fileName)
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

	url := fmt.Sprintf("%s/%s", u.PublicURL, key)
	return url, nil
}
