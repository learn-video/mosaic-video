package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/fs"
	"github.com/mauricioabreu/mosaic-video/internal/locking"
	"github.com/mauricioabreu/mosaic-video/internal/logging"
	"github.com/mauricioabreu/mosaic-video/internal/worker"
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
		fs.Module,
		fx.Invoke(worker.Run),
	)

	app.Run()
}
