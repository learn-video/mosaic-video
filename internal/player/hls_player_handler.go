package player

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type HlsPlayerHandler struct{}

func NewHlsPlayerHandler() *HlsPlayerHandler {
	return &HlsPlayerHandler{}
}

func (hh *HlsPlayerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	filename := vars["file"]

	currentPath, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filepath := fmt.Sprintf("%s/internal/player/html/index.html", currentPath)

	if filename != "" {
		filepath = fmt.Sprintf("%s/internal/player/html/assets/%s", currentPath, filename)
	}

	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	fileInfo, err := file.Stat()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, req, fileInfo.Name(), fileInfo.ModTime(), file)
}
