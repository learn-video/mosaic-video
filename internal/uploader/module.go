package uploader

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mauricioabreu/mosaic-video/internal/config"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewHandler,
)

func NewHandler(cfg *config.Config) *FileUploadHandler {
	return &FileUploadHandler{
		BaseDir: cfg.AssetsPath,
	}
}

func Run(lifecycle fx.Lifecycle, handler *FileUploadHandler) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			r := mux.NewRouter()
			r.Handle("/hls/{folder}/{filename:[a-zA-Z0-9/_-]+}.{ext:[a-zA-Z0-9_-]+}", handler).Methods("PUT", "POST")
			go http.ListenAndServe(":8080", r)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
