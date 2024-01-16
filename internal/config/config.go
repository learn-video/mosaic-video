package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	Redis struct {
		Host string `env:"REDIS_HOST"`
		Port string `env:"REDIS_PORT"`
	}
	API struct {
		URL string `env:"MOSAICS_API_URL"`
	}
	StaticsPath string `env:"STATICS_PATH"`
	S3          struct {
		Endpoint        string `env:"S3_ENDPOINT"`
		AccessKeyID     string `env:"S3_ACCESS_KEY_ID"`
		SecretAccessKey string `env:"S3_SECRET_ACCESS_KEY"`
		BucketName      string `env:"S3_BUCKET_NAME"`
	}
	UploaderEndpoint string `env:"UPLOADER_ENDPOINT"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
