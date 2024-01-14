package mosaic

import (
	"encoding/json"
	"net/http"
)

type (
	Position struct {
		X int
		Y int
	}
	Media struct {
		URL      string `json:"url"`
		Position Position
		Scale    string
	}

	Mosaic struct {
		Name          string  `json:"name"`
		BackgroundUrl string  `json:"background_url"`
		Medias        []Media `json:"medias"`
		WithAudio     bool    `json:"with_audio"`
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
