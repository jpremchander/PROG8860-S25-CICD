package config

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var S3Client *s3.S3
var S3Session *session.Session

func InitS3() {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ca-central-1"
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatal("Failed to create AWS session:", err)
	}

	S3Session = sess
	S3Client = s3.New(sess)
	log.Println("AWS S3 initialized!")
}
