package config

import (
	"errors"
	"os"
)

const (
	minioHostEnvName = "MINIO_HOST"
	minioPortEnvName = "MINIO_PORT"
	minioDSN         = "MINIO_DSN"
	minioAccessKey   = "MINIO_USER_SERVICE_ACCESS_KEY"
	minioSecretKey   = "MINIO_USER_SERVICE_SECRET_KEY"
)

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
}

func NewMinioConfig() (*MinioConfig, error) {
	dsn := os.Getenv(minioDSN)
	if len(dsn) == 0 {
		return nil, errors.New("minio host not found")
	}

	accessKey := os.Getenv(minioAccessKey)
	if len(accessKey) == 0 {
		return nil, errors.New("minio accessKey not found")
	}

	secretKey := os.Getenv(minioSecretKey)
	if len(secretKey) == 0 {
		return nil, errors.New("minio accessKey not found")
	}

	return &MinioConfig{
		Endpoint:  dsn,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}, nil
}
