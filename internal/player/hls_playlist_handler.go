package player

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mauricioabreu/mosaic-video/internal/storage/s3"
)

type HlsPlaylistHandler struct {
	s3Client *s3.Client
}

func NewHlsPlaylistHandler(s3c *s3.Client) *HlsPlaylistHandler {
	return &HlsPlaylistHandler{
		s3Client: s3c,
	}
}

func (hh *HlsPlaylistHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	filename := vars["filename"]
	hh.servePlaylistHTTPImpl(filename, w, req)
}

func (hh *HlsPlaylistHandler) servePlaylistHTTPImpl(filename string, w http.ResponseWriter, req *http.Request) {
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
