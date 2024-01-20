package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type LocalStorage struct {
	Path string `env:"LOCAL_STORAGE_PATH,notEmpty"`
}

type S3 struct {
	Endpoint         string `env:"S3_ENDPOINT,notEmpty"`
	AccessKeyID      string `env:"S3_ACCESS_KEY_ID,notEmpty"`
	SecretAccessKey  string `env:"S3_SECRET_ACCESS_KEY,notEmpty"`
	BucketName       string `env:"S3_BUCKET_NAME,notEmpty"`
	UploaderEndpoint string `env:"UPLOADER_ENDPOINT,notEmpty"`
}

type Config struct {
	Redis struct {
		Host string `env:"REDIS_HOST,notEmpty"`
		Port string `env:"REDIS_PORT,notEmpty"`
	}
	API struct {
		URL string `env:"MOSAICS_API_URL,notEmpty"`
	}
	StorageType  string `env:"STORAGE_TYPE" envDefault:"local"`
	LocalStorage LocalStorage
	S3           S3
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	switch cfg.StorageType {
	case "local":
		if err := env.Parse(&cfg.LocalStorage); err != nil {
			return nil, err
		}
	case "s3":
		if err := env.Parse(&cfg.S3); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid storage type: %s", cfg.StorageType)
	}

	return &cfg, nil
}
