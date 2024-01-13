package player

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewHlsPlaylistHandler,
	NewHlsFragmentHandler,
	NewHlsPlayerHandler,
)

func Run(
	lifecycle fx.Lifecycle,
	playlistHandler *HlsPlaylistHandler,
	fragmentHandler *HlsFragmentHandler,
	playerHandler *HlsPlayerHandler,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			r := mux.NewRouter()
			r.Handle("/playlist/{filename}", playlistHandler).Methods("GET")
			r.Handle("/fragment/{filename}", fragmentHandler).Methods("GET")
			r.Handle("/player", playerHandler).Methods("GET")
			go http.ListenAndServe(":8090", r)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
