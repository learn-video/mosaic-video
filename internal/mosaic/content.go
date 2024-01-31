package mosaic

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
		IsLoop   bool `json:"is_loop"`
	}

	Mosaic struct {
		Name          string  `json:"name"`
		BackgroundURL string  `json:"background_url"`
		Medias        []Media `json:"medias"`
		Audio         Audio   `json:"audio"`
	}
)
