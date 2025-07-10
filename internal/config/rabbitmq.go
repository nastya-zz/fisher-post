package config

import (
	"errors"
	"os"
)

const (
	dsnRMQEnvName = "RMQ_DSN"
)

type RMQConfig interface {
	DSN() string
}

type rmqConfig struct {
	dsn string
}

func NewRMQConfig() (RMQConfig, error) {
	dsn := os.Getenv(dsnRMQEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("rmq dsn not found")
	}

	return &rmqConfig{
		dsn: dsn,
	}, nil
}

func (cfg *rmqConfig) DSN() string {
	return cfg.dsn
}
