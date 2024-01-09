package uploader

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
)

type FileUploadHandler struct {
	BaseDir string
}

func (fu *FileUploadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	curFileURL := req.URL.EscapedPath()[len("/hls"):]
	vars := mux.Vars(req)
	folder := vars["folder"]
	curFolderPath := path.Join(fu.BaseDir, folder)
	curFilePath := path.Join(fu.BaseDir, curFileURL)
	fu.serveHTTPImpl(curFolderPath, curFilePath, w, req)
}

func (fu *FileUploadHandler) serveHTTPImpl(curFolderPath string, curFilePath string, w http.ResponseWriter, req *http.Request) {
	if _, err := os.Stat(curFolderPath); os.IsNotExist(err) {
		err := os.MkdirAll(curFolderPath, os.ModePerm)
		if err != nil {
			log.Printf("fail to create file %v", err)
		}
	}

	if _, err := os.Stat(curFilePath); err == nil {
		log.Printf("rewrite file %s @ %v \n", curFilePath, time.Now().Format(time.RFC3339))
		data, _ := io.ReadAll(req.Body)
		err = os.WriteFile(curFilePath, data, 0644)
		if err != nil {
			log.Printf("fail to create file %v \n", err)
		}
		return
	}

	f, rerr := os.Create(curFilePath)
	if rerr != nil {
		log.Printf("fail to create file %s : %v\n", curFilePath, rerr)
		return
	}
	defer f.Close()

	_, rerr = io.Copy(f, req.Body)
	if rerr != nil {
		log.Printf("fail to create file %v \n", rerr)
	}
}
