package uploader

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewHandler,
)

func Run(lifecycle fx.Lifecycle, handler *FileUploadHandler) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			r := mux.NewRouter()
			r.Handle("/hls/{folder}/{filename}", handler).Methods("PUT", "POST")
			go http.ListenAndServe(":8080", r)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
