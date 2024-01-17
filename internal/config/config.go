package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	Redis struct {
		Host string `env:"REDIS_HOST,notEmpty"`
		Port string `env:"REDIS_PORT,notEmpty"`
	}
	API struct {
		URL string `env:"MOSAICS_API_URL,notEmpty"`
	}
	StaticsPath string `env:"STATICS_PATH,notEmpty"`
	S3          struct {
		Endpoint        string `env:"S3_ENDPOINT,notEmpty"`
		AccessKeyID     string `env:"S3_ACCESS_KEY_ID,notEmpty"`
		SecretAccessKey string `env:"S3_SECRET_ACCESS_KEY,notEmpty"`
		BucketName      string `env:"S3_BUCKET_NAME,notEmpty"`
	}
	UploaderEndpoint string `env:"UPLOADER_ENDPOINT,notEmpty"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
