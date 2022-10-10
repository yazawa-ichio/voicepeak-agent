package voicepeakagent

import (
	"fmt"
	"path/filepath"
	"strings"
)

type SayRequest struct {
	Text     string          `json:"text"`
	FileName string          `json:"file_name"`
	Narrator Narrator        `json:"narrator,omitempty"`
	Speed    int             `json:"speed,omitempty"`
	Pitch    int             `json:"pitch,omitempty"`
	Emotions []*EmotionValue `json:"emotions,omitempty"`
}

type EmotionValue struct {
	Emotion Emotion `json:"emotion"`
	Rate    int     `json:"rate"`
}

func NewSayRequest(text, fileName string) *SayRequest {
	return &SayRequest{
		Text:     text,
		FileName: fileName,
	}
}

func (s *SayRequest) SetEmotion(emotion Emotion, rate int) {
	s.Emotions = append(s.Emotions, &EmotionValue{
		Emotion: emotion,
		Rate:    rate,
	})
}

func (sr *SayRequest) GetArgs(outDir string) ([]string, error) {
	var args []string
	args = append(args, "-s")
	args = append(args, sr.Text)
	args = append(args, "-o")
	args = append(args, filepath.Join(outDir, sr.FileName))
	if len(sr.Narrator) > 0 {
		args = append(args, "-n")
		args = append(args, sr.Narrator.String())
	}
	if sr.Pitch != 0 {
		pitch := sr.Pitch
		if pitch < -300 || 300 < pitch {
			return nil, fmt.Errorf("pitch (-300 - 300) value:%v", pitch)
		}
		args = append(args, "--pitch")
		args = append(args, fmt.Sprintf("%v", pitch))
	}
	if sr.Speed != 0 {
		speed := sr.Speed
		if speed < 50 || 200 < speed {
			return nil, fmt.Errorf("speed (50 - 200) value:%v", speed)
		}
		args = append(args, "--speed")
		args = append(args, fmt.Sprintf("%v", speed))
	}
	if len(sr.Emotions) > 0 {
		args = append(args, "-e")
		args = append(args, sr.getEmotionArg())
	}
	return args, nil
}

func (sr *SayRequest) getEmotionArg() string {
	var list []string
	for _, e := range sr.Emotions {
		list = append(list, fmt.Sprintf("%s=%v", e.Emotion, e.Rate))
	}
	return strings.Join(list, ",")
}
