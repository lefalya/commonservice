package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
)

func GeneratePresignURL(key string, fileSize int64, ttl int64) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("loading config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	presignInput := &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("MIXPHOTO_BUCKET")),
		Key:    aws.String(key),
	}

	presignClient := s3.NewPresignClient(client)
	presignResult, err := presignClient.PresignPutObject(context.TODO(), presignInput,
		func(opts *s3.PresignOptions) {
			// TODO ini masih belum bisa kalau di kasih ttl
			// opts.Expires = time.Duration(ttl)
		},
	)
	if err != nil {
		return "", fmt.Errorf("presigning URL: %w", err)
	}

	return presignResult.URL, nil
}

func GetObject(key string) error {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	file, err := os.Create(key)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	downloader := manager.NewDownloader(client)

	_, err = downloader.Download(context.TODO(), file, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("MIXPHOTO_BUCKET")),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("downloading file: %w", err)
	}

	fmt.Printf("File downloaded successfully to %s\n", key)
	return nil
}

func UploadObject(filePath, key string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	// Create an uploader
	uploader := manager.NewUploader(client)

	// Upload the file
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("MIXPHOTO_BUCKET")),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("uploading file: %w", err)
	}

	fmt.Println("File uploaded successfully to S3 with key:", key)
	return nil
}
