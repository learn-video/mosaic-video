package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mauricioabreu/mosaic-video/config"
	"github.com/mauricioabreu/mosaic-video/locking"
	"github.com/mauricioabreu/mosaic-video/logging"
	"github.com/mauricioabreu/mosaic-video/worker"
	"go.uber.org/fx"
)

func main() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	if err := godotenv.Load(envFile); err != nil {
		log.Println("Could not load .env file")
	}

	app := fx.New(
		config.Module,
		fx.Provide(
			logging.NewLogger,
			locking.NewRedisLocker,
		),
		fx.Invoke(worker.Run),
	)

	app.Run()
}
