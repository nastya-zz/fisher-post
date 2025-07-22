package minio

import (
	"context"
	"fmt"
	"strings"

	"github.com/minio/minio-go/v7"

	"post/pkg/logger"
)

func (s Service) RemoveFile(ctx context.Context, link string) error {
	logger.Warn("service.minio.RemoveFile", "link", link)
	const op = "service.minio.RemoveFile"
	fileName, err := s.getFileName(link)
	if err != nil {
		return err
	}

	err = s.minioClient.Client.RemoveObject(ctx, s.minioClient.BucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("не удалось удалить объект: %v", err)
	}

	logger.Info(op, "Файл успешно удален", "бакет", s.minioClient.BucketName, "имя файла", fileName)
	return nil
}

func (s Service) getFileName(link string) (string, error) {
	name, found := strings.CutPrefix(link, s.minioClient.BucketName+"/")
	if !found {
		return "", fmt.Errorf("cannot get filename")
	}

	return name, nil
}
