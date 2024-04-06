package s3

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client *s3.Client
}

type S3Config struct {
	AccessKey    string
	SecretKey    string
	Region       string
	SessionToken string
}

func GetClient() *S3Client {
	s3Config := GetS3Config()
	opts := s3.Options{
		Region: *aws.String(s3Config.Region),
		// Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(s3Config.AccessKey, s3Config.SecretKey, s3Config.SessionToken)),
	}

	// Create an Amazon S3 service client
	s3Client := s3.New(opts)

	client := &S3Client{
		Client: s3Client,
	}

	return client
}

func GetS3Config() *S3Config {
	return &S3Config{
		AccessKey:    os.Getenv("S3_ACCESS_KEY"),
		SecretKey:    os.Getenv("S3_SECRET_KEY"),
		Region:       os.Getenv("S3_REGION"),
		SessionToken: os.Getenv("S3_SESSION_TOKEN"),
	}
}
