package player

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
)

type HlsManifestHandler struct {
	cfg *config.Config
}

func NewHlsManifestHandler(cfg *config.Config) *HlsManifestHandler {
	return &HlsManifestHandler{
		cfg: cfg,
	}
}

func (h *HlsManifestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	tasks, err := mosaic.FetchMosaicTasks(h.cfg.API.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	manifest := make([]Manifest, len(tasks))
	for i, m := range tasks {
		manifest[i].Name = m.Name
		manifest[i].PlaylistURL = fmt.Sprintf("%s/playlist/playlist-%s.m3u8", h.cfg.PlayerEndpoint, m.Name)
	}

	output, err := json.Marshal(&manifest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(output)
}
