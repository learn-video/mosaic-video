package player

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mauricioabreu/mosaic-video/internal/storage/s3"
)

type HlsFragmentHandler struct {
	s3Client *s3.Client
}

func NewHlsFragmentHandler(s3c *s3.Client) *HlsFragmentHandler {
	return &HlsFragmentHandler{
		s3Client: s3c,
	}
}

func (hh *HlsFragmentHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	filename := vars["filename"]
	hh.serveFragmentHTTPImpl(filename, w, req)
}

func (hh *HlsFragmentHandler) serveFragmentHTTPImpl(filename string, w http.ResponseWriter, req *http.Request) {
	content, err := hh.s3Client.Get(filename)
	if err != nil {
		log.Printf("failure to get %s file from bucket, err: %v", filename, err)
		w.WriteHeader(http.StatusBadRequest)
	}

	if _, err := io.Copy(w, content); err != nil {
		log.Printf("failure copy %s file to response, err: %v", filename, err)
		w.WriteHeader(http.StatusBadRequest)
	}
}
