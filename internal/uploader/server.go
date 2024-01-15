package uploader

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mauricioabreu/mosaic-video/internal/storage"
	"go.uber.org/zap"
)

type FileUploadHandler struct {
	storageHandler storage.Storage
	logger         *zap.SugaredLogger
}

func NewHandler(s storage.Storage, logger *zap.SugaredLogger) *FileUploadHandler {
	return &FileUploadHandler{
		storageHandler: s,
		logger:         logger,
	}
}

func (fu *FileUploadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	folder := vars["folder"]
	filename := vars["filename"]
	fu.serveHTTPImpl(folder, filename, w, req)
}

func (fu *FileUploadHandler) serveHTTPImpl(folder, filename string, w http.ResponseWriter, req *http.Request) {
	fu.logger.Debugf("uploading file %s to folder %s", filename, folder)

	data, err := io.ReadAll(req.Body)
	if err != nil {
		fu.logger.Errorf("failed to read request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	filepath := fmt.Sprintf("%s/%s", folder, filename)

	if err := fu.storageHandler.Upload(filepath, data); err != nil {
		fu.logger.Errorf("failed to upload file to storage: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	fu.logger.Infof("file %s uploaded to folder %s", filename, folder)
}
