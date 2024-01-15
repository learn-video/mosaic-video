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
	return a == Audio(NoAudio)
}

func (a Audio) IsFirstInput() bool {
	return a == Audio(FirstInput)
}

func (a Audio) IsAllInputs() bool {
	return a == Audio(AllInputs)
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
