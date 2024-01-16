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
		BackgroundURL string  `json:"background_url"`
		Medias        []Media `json:"medias"`
		WithAudio     bool    `json:"with_audio"`
	}
)

func FetchMosaicTasks(apiURL string) ([]Mosaic, error) {
	//nolint:gosec // we are skipping this because it's better way to validate the application, for while.
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var tasks []Mosaic

	err = json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		return nil, err
	}

	return tasks, err
}
