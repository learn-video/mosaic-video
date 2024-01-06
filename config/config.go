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
	AssetsPath string
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
		AssetsPath: os.Getenv("ASSETS_PATH"),
	}
}
