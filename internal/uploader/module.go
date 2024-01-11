package uploader

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mauricioabreu/mosaic-video/internal/storage/s3"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewHandler,
)

func NewHandler(s3c *s3.Client) *FileUploadHandler {
	return &FileUploadHandler{
		s3Client: s3c,
	}
}

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
