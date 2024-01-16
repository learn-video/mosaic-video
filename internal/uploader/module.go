package uploader

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

			go func() {
				server := &http.Server{
					Addr:              ":8080",
					Handler:           r,
					ReadHeaderTimeout: time.Duration(0),
				}

				if err := server.ListenAndServe(); err != nil {
					panic(fmt.Errorf("failed to start uploader on :8080 port. err=%v", err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
