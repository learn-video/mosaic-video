package config

import "os"

type Config struct {
	Redis struct {
		Host string
		Port string
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
	}
}
