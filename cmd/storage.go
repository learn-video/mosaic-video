package cmd

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/logging"
	"github.com/mauricioabreu/mosaic-video/internal/storage"
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
					fx.Annotate(
						s3.NewClient,
						fx.As(new(storage.Storage)),
					),
					uploader.NewHandler,
				),
				fx.Invoke(uploader.Run),
			)

			app.Run()
		},
	}
}
