package config

import "os"

type Config struct {
	Redis struct {
		Host string
		Port string
	}
	API struct {
		URL string
	}
	StaticsPath string
	S3          struct {
		Endpoint        string
		AccessKeyID     string
		SecretAccessKey string
		BucketName      string
	}
}

func NewConfig() *Config {
	return &Config{
		Redis: struct {
			Host string
			Port string
		}{
			Host: os.Getenv("REDIS_HOST"),
			Port: os.Getenv("REDIS_PORT"),
		},
		API: struct {
			URL string
		}{
			URL: os.Getenv("MOSAICS_API_URL"),
		},
		S3: struct {
			Endpoint        string
			AccessKeyID     string
			SecretAccessKey string
			BucketName      string
		}{
			Endpoint:        os.Getenv("S3_ENDPOINT"),
			AccessKeyID:     os.Getenv("S3_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
			BucketName:      os.Getenv("S3_BUCKET_NAME"),
		},
		StaticsPath: os.Getenv("STATICS_PATH"),
	}
}
