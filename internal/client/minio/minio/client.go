package minio

import (
	"context"
	"fmt"
	miniolib "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	Client     *miniolib.Client
	BucketName string
	endpoint   string
}

const bucketName = "user-service"

func New(ctx context.Context, endpoint, accessKeyID, secretAccessKey string) (*Client, error) {

	minioClient, err := miniolib.New(endpoint, &miniolib.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		return nil, fmt.Errorf("произошла ошибка при  создании клиента минио %w", err)
	}

	exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
	if errBucketExists != nil && !exists {
		return nil, fmt.Errorf("произошла ошибка при  проверке бакета %s %w", bucketName, err)
	}

	return &Client{
		Client:     minioClient,
		BucketName: bucketName,
		endpoint:   endpoint,
	}, err
}

func (m Client) SourcePrefix() string {
	return fmt.Sprintf("%s/", m.BucketName)
}
