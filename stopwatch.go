package voicepeakagent

import (
	"time"
)

type stopwatch struct {
	tag       string
	startTime time.Time
}

func newStopwatch(tag string) *stopwatch {
	return &stopwatch{
		tag:       tag,
		startTime: time.Now(),
	}
}

func (s *stopwatch) Stop() {
	TraceLog("%s %s", s.tag, time.Since(s.startTime))
}
