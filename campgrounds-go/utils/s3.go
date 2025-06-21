package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"campgrounds-go/config"
)

func UploadToS3(file *multipart.FileHeader) (string, string, error) {
	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	// Read file content
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, src); err != nil {
		return "", "", err
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	key := fmt.Sprintf("campgrounds/%s", filename)

	bucketName := os.Getenv("AWS_S3_BUCKET")
	if bucketName == "" {
		bucketName = "yelp-camp-project"
	}

	// Upload to S3
	_, err = config.S3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String(getContentType(file.Filename)),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		return "", "", err
	}

	// Generate URL
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", 
		bucketName, 
		os.Getenv("AWS_REGION"), 
		key)

	return url, filename, nil
}

func getContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}
