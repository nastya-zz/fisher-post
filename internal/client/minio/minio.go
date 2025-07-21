package minio

import (
	"github.com/minio/minio-go/v7"
)

type ClientMinio interface {
	Connect() *minio.Client
}
