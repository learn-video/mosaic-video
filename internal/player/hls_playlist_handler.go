package player

import (
	"fmt"
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
	folder := vars["folder"]
	filename := vars["filename"]
	hh.servePlaylistHTTPImpl(folder, filename, w, req)
}

func (hh *HlsPlaylistHandler) servePlaylistHTTPImpl(folder string, filename string, w http.ResponseWriter, req *http.Request) {
	file := fmt.Sprintf("%s/%s", folder, filename)
	content, err := hh.s3Client.Get(file)
	if err != nil {
		log.Printf("failure to get %s file from bucket, err: %v", file, err)
		w.WriteHeader(http.StatusBadRequest)
	}

	if _, err := io.Copy(w, content); err != nil {
		log.Printf("failure copy %s file to response, err: %v", file, err)
		w.WriteHeader(http.StatusBadRequest)
	}
}
