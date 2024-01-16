package mosaic

import (
	"encoding/json"
	"net/http"
)

type Audio string

const (
	NoAudio    Audio = "no_audio"
	FirstInput Audio = "first_input"
	AllInputs  Audio = "all_inputs"
)

func (a Audio) IsNoAudio() bool {
	return a == NoAudio
}

func (a Audio) IsFirstInput() bool {
	return a == FirstInput
}

func (a Audio) IsAllInputs() bool {
	return a == AllInputs
}

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
		Audio         Audio   `json:"audio"`
	}
)

func FetchMosaicTasks(apiURL string) ([]Mosaic, error) {
	//nolint:gosec // we are skipping this because it's better way to validate the application, for now.
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
