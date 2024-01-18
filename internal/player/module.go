package player

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewHlsPlaylistHandler,
	NewHlsPlayerHandler,
)

func Run(
	lifecycle fx.Lifecycle,
	playlistHandler *HlsPlaylistHandler,
	playerHandler *HlsPlayerHandler,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			r := mux.NewRouter()
			r.Handle("/playlist/{folder}/{filename}", playlistHandler).Methods("GET")
			r.Handle("/player", playerHandler).Methods("GET")
			r.Handle("/player/assets/{file}", playerHandler).Methods("GET")

			go func() {
				server := &http.Server{
					Addr:              ":8090",
					Handler:           r,
					ReadHeaderTimeout: time.Duration(0),
				}

				if err := server.ListenAndServe(); err != nil {
					panic(fmt.Errorf("failed to start player on :8090 port. err=%v", err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
