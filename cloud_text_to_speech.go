package pollydent

import (
	"io"
)

type Request struct {
	Input       input       `json:"input"`
	Voice       voice       `json:"voice"`
	AudioConfig audioConfig `json:"audioConfig"`
}

type input struct {
	// Text string `json:"text"`
	SSML string `json:"ssml"`
}

type voice struct {
	LanguageCode string `json:"languageCode"`
	Name         string `json:"name"`
	SsmlGender   string `json:"ssmlGender"`
}

type audioConfig struct {
	AudioEncoding   string `json:"audioEncoding"`
	SampleRateHertz int    `json:"sampleRateHertz"`
}

type Response struct {
	AudioContent string `json:"audioContent"`
}

type MP3Wrapper struct {
	io.Reader
}

func (w *MP3Wrapper) Close() error {
	return nil
}
