package minio

import (
	"post/internal/client/minio/minio"
)

type Service struct {
	minioClient *minio.Client
}

func New(minioClient *minio.Client) *Service {
	return &Service{
		minioClient: minioClient,
	}
}
