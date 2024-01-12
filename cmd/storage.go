package cmd

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/storage"
	"github.com/mauricioabreu/mosaic-video/internal/uploader"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func Store() *cobra.Command {
	return &cobra.Command{
		Use:   "storage",
		Short: "Start backend to store mosaics",
		Run: func(cmd *cobra.Command, args []string) {
			if err := godotenv.Load(".env", ".penv"); err != nil {
				log.Println("Could not load .env file")
			}

			app := fx.New(
				config.Module,
				uploader.Module,
				storage.Module,
				fx.Invoke(uploader.Run),
			)

			app.Run()
		},
	}
}
