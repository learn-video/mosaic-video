package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/locking"
	"github.com/mauricioabreu/mosaic-video/internal/logging"
	"github.com/mauricioabreu/mosaic-video/internal/uploader"
	"github.com/mauricioabreu/mosaic-video/internal/worker"
	"go.uber.org/fx"
)

func main() {
	if err := godotenv.Load(".env", ".penv"); err != nil {
		log.Println("Could not load .env file")
	}

	app := fx.New(
		config.Module,
		uploader.Module,
		fx.Provide(
			logging.NewLogger,
			locking.NewRedisLocker,
		),
		fx.Invoke(uploader.Run),
		fx.Invoke(worker.Run),
	)

	app.Run()
}
