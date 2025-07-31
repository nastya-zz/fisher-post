package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	dsnRedisEnvName      = "REDIS_DSN"
	passwordRedisEnvName = "REDIS_PASSWORD"
	dbRedisEnvName       = "REDIS_DB"
)

type RedisConfig interface {
	Config() *redisConfig
}

type redisConfig struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisConfig() (RedisConfig, error) {
	dsn := os.Getenv(dsnRedisEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("redis dsn not found")
	}

	password := os.Getenv(passwordRedisEnvName)
	if len(password) == 0 {
		return nil, errors.New("redis password not found")
	}

	db, err := strconv.Atoi(os.Getenv(dbRedisEnvName))
	if err != nil {
		return nil, errors.New("redis db not found")
	}

	return &redisConfig{
		Addr:     dsn,
		Password: password,
		DB:       db,
	}, nil
}

func (cfg *redisConfig) Config() *redisConfig {
	return cfg
}
