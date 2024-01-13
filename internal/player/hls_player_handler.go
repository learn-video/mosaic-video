package player

import (
	"fmt"
	"net/http"
	"os"
)

type HlsPlayerHandler struct{}

func NewHlsPlayerHandler() *HlsPlayerHandler {
	return &HlsPlayerHandler{}
}

func (hh *HlsPlayerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	currentPath, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	indexPath := fmt.Sprintf("%s/internal/player/html/index.html", currentPath)
	file, err := os.Open(indexPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.ServeContent(w, req, fileInfo.Name(), fileInfo.ModTime(), file)
	return
}
