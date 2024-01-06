package mosaic

import (
	"encoding/json"
	"net/http"
)

type (
	Media struct {
		URL string `json:"url"`
	}

	Mosaic struct {
		Name      string  `json:"name"`
		Medias    []Media `json:"medias"`
		WithAudio bool    `json:"with_audio"`
	}
)

func FetchMosaicTasks(apiURL string) ([]Mosaic, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var tasks []Mosaic
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}

	return tasks, err
}
