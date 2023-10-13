package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mauricioabreu/mosaic-video/config"
	"github.com/mauricioabreu/mosaic-video/locking"
	"github.com/mauricioabreu/mosaic-video/logging"
	"github.com/mauricioabreu/mosaic-video/worker"
	"go.uber.org/fx"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
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
