package pollydent

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/hajimehoshi/oto"
)

type SpeechParams struct {
	Message string
	Voice   string
	Speed   int
}

type Speaker interface {
	Send(SpeechParams) (io.Reader, error)
}

type PollySpeaker struct {
	config *PollyConfig
	sess   *session.Session
}

func (p *PollySpeaker) Send(config SpeechParams) (io.Reader, error) {
	var err error

	if config.Speed == 0 {
		config.Speed = p.config.Speed
	}

	if config.Voice == "" {
		config.Voice = p.config.Voice
	}

	text := `<speak><prosody rate="` + strconv.Itoa(config.Speed) + `%"><![CDATA[` + config.Message + `]]></prosody></speak>`

	pol := polly.New(p.sess, aws.NewConfig().WithRegion(p.config.Region))

	params := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String(p.config.Format),
		Text:         aws.String(text),
		TextType:     aws.String(p.config.TextType),
		VoiceId:      aws.String(config.Voice),
	}

	resp, err := pol.SynthesizeSpeech(params)
	if err != nil {
		return nil, err
	}
	return resp.AudioStream, nil
}

// Pollydent is structure to manage read aloud
type Pollydent struct {
	playMutex   *sync.Mutex
	audioConfig AudioConfig
	speaker     Speaker
}

// NewPollydent news Polly structure
func NewPollydentWithPolly(accessKey, secretKey string, config *PollyConfig) (*Pollydent, error) {
	if accessKey == "" || secretKey == "" {
		return nil, errors.New("Access key or Secret key are invalid")
	}

	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	sess := session.New(&aws.Config{Credentials: creds})

	if config == nil {
		config = defaultConfig()
	}

	return &Pollydent{
		playMutex:   new(sync.Mutex),
		audioConfig: &PollyAudioConfig{},
		speaker:     &PollySpeaker{config, sess},
	}, nil
}

func (p *Pollydent) Play(reader io.Reader) (err error) {
	p.playMutex.Lock()
	defer p.playMutex.Unlock()

	totalData := make([]byte, 0)
	for {
		var n int
		data := make([]byte, 65535)
		if n, err = reader.Read(data); err != nil {
			if err != io.EOF {
				return
			}
			totalData = append(totalData, data[:n]...)
			break
		} else {
			totalData = append(totalData, data[:n]...)
		}
	}

	playerCtx, err := oto.NewContext(
		p.audioConfig.SampleRate(),
		p.audioConfig.NumOfChanel(),
		p.audioConfig.ByteParSample(),
		len(totalData))
	if err != nil {
		return
	}
	defer playerCtx.Close()

	player := playerCtx.NewPlayer()

	defer player.Close()

	if _, err = player.Write(totalData); err != nil {
		return
	}

	return
}

func (p *Pollydent) SendToServer(param SpeechParams) (io.Reader, error) {
	return p.speaker.Send(param)
}

// ReadAloud reads aloud msg by Polly
func (p *Pollydent) ReadAloud(msg string) (err error) {
	if msgLen := len([]rune(msg)); msgLen > 1500 {
		const errMsg = "message size is %d. Please pass with the length of 1500 or less"
		err = fmt.Errorf(errMsg, msgLen)
		return err
	}

	reader, err := p.speaker.Send(SpeechParams{Message: msg})
	if err != nil {
		return
	}
	p.Play(reader)
	return
}
