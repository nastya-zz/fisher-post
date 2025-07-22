package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"path/filepath"
)

func (s Service) UploadFile(ctx context.Context, file []byte, filename string) (string, error) {
	const op = "service.UploadFile"

	ext := filepath.Ext(filename)
	contentType, err := getContentTYpe(ext)
	if err != nil {
		return "", err
	}

	ui, err := s.minioClient.Client.PutObject(ctx, s.minioClient.BucketName, filename, bytes.NewReader(file), int64(len(file)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	srcLink := s.minioClient.SourcePrefix()
	link := fmt.Sprintf("%s%s", srcLink, filename)

	fmt.Printf("%s: uploaded file %s to bucket %s\n", op, ui.ETag, s.minioClient.BucketName)

	return link, nil
}

func getContentTYpe(ext string) (string, error) {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg", nil
	case ".png":
		return "image/png", nil
	case ".webp":
		return "image/webp", nil
	case ".heic":
		return "image/heic", nil

	default:
		return "", fmt.Errorf("extention not support")
	}
}
