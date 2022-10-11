package speaker

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	voicepeakagent "github.com/yazawa-ichio/voicepeak-agent"
)

type Speaker struct {
	app  *voicepeakagent.App
	Gain float64
}

func NewSpeaker(app *voicepeakagent.App) *Speaker {
	return &Speaker{app: app, Gain: 0}
}

func (s *Speaker) Play(req *voicepeakagent.SayRequest) error {
	return s.PlayWithGain(req, 0)
}

func (s *Speaker) PlayWithGain(req *voicepeakagent.SayRequest, gain float64) error {
	return s.app.Say(req, func(path string) error {

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		streamer, format, err := wav.Decode(f)
		if err != nil {
			return err
		}
		defer streamer.Close()
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		done := make(chan bool)
		speaker.Play(beep.Seq(&effects.Gain{Streamer: streamer, Gain: s.Gain + gain}, beep.Callback(func() {
			done <- true
		})))
		<-done
		return nil
	})
}
