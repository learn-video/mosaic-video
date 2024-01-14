package cmd

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/logging"
	"github.com/mauricioabreu/mosaic-video/internal/player"
	"github.com/mauricioabreu/mosaic-video/internal/storage"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func Player() *cobra.Command {
	return &cobra.Command{
		Use:   "player",
		Short: "Start player to serve hls content",
		Run: func(cmd *cobra.Command, args []string) {
			if err := godotenv.Load(".env", ".penv"); err != nil {
				log.Println("Could not load .env file")
			}

			app := fx.New(
				config.Module,
				storage.Module,
				fx.Provide(
					logging.NewLogger,
					player.NewHlsPlaylistHandler,
					player.NewHlsPlayerHandler,
				),
				fx.Invoke(player.Run),
			)

			app.Run()
		},
	}
}
