package uploader

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mauricioabreu/mosaic-video/internal/storage/s3"
)

type FileUploadHandler struct {
	s3Client *s3.Client
}

func (fu *FileUploadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.EscapedPath()[len("/hls"):]
	vars := mux.Vars(req)
	folder := vars["folder"]
	fu.serveHTTPImpl(folder, filename, w, req)
}

func (fu *FileUploadHandler) serveHTTPImpl(folder, filename string, w http.ResponseWriter, req *http.Request) {
	log.Printf("uploading file %s to folder %s", filename, folder)

	data, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("failed to read request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := fu.s3Client.Upload(filename, data); err != nil {
		log.Printf("failed to upload file to storage: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
