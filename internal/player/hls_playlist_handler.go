package player

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/storage"
	"go.uber.org/zap"
)

type HlsPlaylistHandler struct {
	cfg            *config.Config
	storageHandler storage.Storage
	logger         *zap.SugaredLogger
}

func NewHlsPlaylistHandler(
	cfg *config.Config,
	s storage.Storage,
	logger *zap.SugaredLogger,
) *HlsPlaylistHandler {
	return &HlsPlaylistHandler{
		cfg:            cfg,
		storageHandler: s,
		logger:         logger,
	}
}

func (hh *HlsPlaylistHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	folder := vars["folder"]
	filename := vars["filename"]
	hh.servePlaylistHTTPImpl(folder, filename, w, req)
}

func (hh *HlsPlaylistHandler) servePlaylistHTTPImpl(folder, filename string, w http.ResponseWriter, _ *http.Request) {
	file, err := hh.getFile(folder, filename)

	if err != nil {
		hh.logger.Errorf("failed to get %s file, err: %v", filename, err)

		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		return
	}

	if _, err := io.Copy(w, file); err != nil {
		hh.logger.Errorf("failed to response %s file, err: %v", file, err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (hh *HlsPlaylistHandler) getFile(folder, name string) (io.Reader, error) {
	filename := fmt.Sprintf("%s/%s", folder, name)

	if hh.cfg.StorageType == "local" {
		path := fmt.Sprintf("%s/%s", hh.cfg.LocalStorage.Path, filename)

		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		return file, nil
	}

	file, err := hh.storageHandler.Get(filename)
	if err != nil {
		return nil, err
	}

	return file, nil
}
