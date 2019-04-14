package pollydent

type AudioConfig interface {
	SampleRate() int
	NumOfChanel() int
	ByteParSample() int
}

type PollyAudioConfig struct {
}

func (c *PollyAudioConfig) SampleRate() int {
	return 16000
}

func (c *PollyAudioConfig) NumOfChanel() int {
	return 1
}

func (c *PollyAudioConfig) ByteParSample() int {
	return 2
}
