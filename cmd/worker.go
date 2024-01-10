package cmd

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/locking"
	"github.com/mauricioabreu/mosaic-video/internal/logging"
	"github.com/mauricioabreu/mosaic-video/internal/uploader"
	"github.com/mauricioabreu/mosaic-video/internal/worker"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func Work() *cobra.Command {
	return &cobra.Command{
		Use:   "worker",
		Short: "Start workers to generate mosaics",
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}
}
