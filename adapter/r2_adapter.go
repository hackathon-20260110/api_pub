package adapter

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hackathon-20260110/api/utils"
)

const R2BucketName = "hackathon-20260110"

type R2Adapter interface {
	UploadImage(image []byte, path string, contentType string) (string, error)
}

func NewR2Adapter(s3Client *s3.Client) R2Adapter {
	return &r2Adapter{client: s3Client}
}

type r2Adapter struct {
	client *s3.Client
}

func (a *r2Adapter) UploadImage(image []byte, path string, contentType string) (string, error) {
	objectKey, err := NormalizeObjectKey(path)
	if err != nil {
		return "", utils.WrapError(err)
	}

	if contentType == "" {
		contentType = http.DetectContentType(image)
	}

	ctx := context.Background()
	_, err = a.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(R2BucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(image),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", utils.WrapError(err)
	}

	return BuildPublicURL(objectKey), nil
}

// NormalizeObjectKey normalizes the object key by removing leading slashes
// and validating the path.
func NormalizeObjectKey(path string) (string, error) {
	key := strings.TrimLeft(path, "/")

	if key == "" {
		return "", errors.New("object key cannot be empty")
	}

	if strings.Contains(key, "..") {
		return "", errors.New("object key cannot contain '..'")
	}

	return key, nil
}

// BuildPublicURL constructs the public URL for an object key using the
// R2_PUBLIC_BASE_URL environment variable.
func BuildPublicURL(objectKey string) string {
	baseURL := os.Getenv("R2_PUBLIC_BASE_URL")
	return baseURL + "/" + objectKey
}

// NewR2ClientFromEnv creates a new S3 client for R2 using environment variables.
func NewR2ClientFromEnv() (*s3.Client, error) {
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	endpoint := os.Getenv("R2_ENDPOINT")

	client := s3.New(s3.Options{
		Region: "auto",
		Credentials: credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			"",
		),
		BaseEndpoint: aws.String(endpoint),
	})

	return client, nil
}
