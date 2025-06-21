package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Service struct {
	client *s3.S3
	bucket string
}

func NewS3Service() (*S3Service, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		return nil, err
	}

	return &S3Service{
		client: s3.New(sess),
		bucket: os.Getenv("S3_BUCKET"),
	}, nil
}

func (s *S3Service) UploadImage(file io.Reader, filename string, userID string) (*UploadResult, error) {
	// Read file content
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, file); err != nil {
		return nil, err
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	ext := filepath.Ext(filename)
	key := fmt.Sprintf("campgrounds/%s/%d%s", userID, timestamp, ext)

	// Upload to S3
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String("image/jpeg"), // Default content type
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		return nil, err
	}

	// Return public URL
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", 
		s.bucket, os.Getenv("AWS_REGION"), key)

	return &UploadResult{
		SecureURL: url,
		PublicID:  key,
	}, nil
}

func (s *S3Service) DeleteImage(key string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

// UploadResult mimics Cloudinary's response structure
type UploadResult struct {
	SecureURL string `json:"secure_url"`
	PublicID  string `json:"public_id"`
}
