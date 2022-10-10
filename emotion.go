package voicepeakagent

type Emotion string

const (
	Happy = Emotion("happy")
	Fun   = Emotion("fun")
	Angry = Emotion("angry")
	Sad   = Emotion("sad")
)

func (n *Emotion) String() string {
	return string(*n)
}
