package uploader

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mauricioabreu/mosaic-video/internal/storage/s3"
	"go.uber.org/zap"
)

type FileUploadHandler struct {
	s3Client *s3.Client
	logger   *zap.SugaredLogger
}

func NewHandler(s3c *s3.Client, logger *zap.SugaredLogger) *FileUploadHandler {
	return &FileUploadHandler{
		s3Client: s3c,
		logger:   logger,
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

	if err := fu.s3Client.Upload(filename, data); err != nil {
		fu.logger.Errorf("failed to upload file to storage: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	fu.logger.Infof("file %s uploaded to folder %s", filename, folder)
}
