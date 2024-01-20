package cmd

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/logging"
	"github.com/mauricioabreu/mosaic-video/internal/storage"
	"github.com/mauricioabreu/mosaic-video/internal/storage/local"
	"github.com/mauricioabreu/mosaic-video/internal/storage/s3"
	"github.com/mauricioabreu/mosaic-video/internal/uploader"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func Store() *cobra.Command {
	return &cobra.Command{
		Use:   "storage",
		Short: "Start backend to store mosaics",
		Run: func(cmd *cobra.Command, args []string) {
			if err := godotenv.Load(".env"); err != nil {
				log.Println("Could not load .env file")
			}

			app := fx.New(
				fx.Provide(
					config.NewConfig,
					logging.NewLogger,
					uploader.NewHandler,
					func(cfg *config.Config) (storage.Storage, error) {
						if cfg.StorageType.IsLocal() {
							return local.NewClient(cfg), nil
						}
						return s3.NewClient(cfg)
					},
				),
				fx.Invoke(uploader.Run),
			)

			app.Run()
		},
	}
}
