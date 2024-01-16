package player

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mauricioabreu/mosaic-video/internal/storage"
	"go.uber.org/zap"
)

type HlsPlaylistHandler struct {
	storageHandler storage.Storage
	logger         *zap.SugaredLogger
}

func NewHlsPlaylistHandler(s storage.Storage, logger *zap.SugaredLogger) *HlsPlaylistHandler {
	return &HlsPlaylistHandler{
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

func (hh *HlsPlaylistHandler) servePlaylistHTTPImpl(folder string, filename string, w http.ResponseWriter, req *http.Request) {
	file := fmt.Sprintf("%s/%s", folder, filename)
	content, err := hh.storageHandler.Get(file)

	if err != nil {
		hh.logger.Errorf("failed to get %s file from bucket, err: %v", file, err)
		w.WriteHeader(http.StatusBadRequest)
	}

	if _, err := io.Copy(w, content); err != nil {
		hh.logger.Errorf("failed to copy %s file to response, err: %v", file, err)
		w.WriteHeader(http.StatusBadRequest)
	}
}
