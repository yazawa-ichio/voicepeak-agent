package main

import (
	"strconv"
	"strings"

	voicepeakagent "github.com/yazawa-ichio/voicepeak-agent"
)

type speakRequest struct {
	Text     string                  `schema:"text,required"`
	Narrator voicepeakagent.Narrator `schema:"narrator"`
	Speed    int                     `schema:"speed"`
	Pitch    int                     `schema:"pitch"`
	Emotions string                  `schema:"emotions"`
	Gain     float64                 `schema:"gain"`
}

func (s *speakRequest) ConvertSayRequest() *voicepeakagent.SayRequest {
	return &voicepeakagent.SayRequest{
		Text:     s.Text,
		Narrator: s.Narrator,
		Speed:    s.Speed,
		Pitch:    s.Pitch,
		Emotions: s.getEmotionValues(),
	}
}

func (s *speakRequest) getEmotionValues() []voicepeakagent.EmotionValue {
	if len(s.Emotions) > 0 {
		var list []voicepeakagent.EmotionValue
		for _, value := range strings.Split(s.Emotions, ",") {
			arr := strings.Split(value, "=")
			rate, _ := strconv.Atoi(arr[1])
			list = append(list, voicepeakagent.EmotionValue{
				Emotion: voicepeakagent.Emotion(arr[0]),
				Rate:    rate,
			})
		}
		return list
	}
	return make([]voicepeakagent.EmotionValue, 0)
}
