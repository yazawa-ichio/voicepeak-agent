package voicepeakagent

import "log"

type LogLevel int

const (
	InfoLogLevel  = LogLevel(0)
	DebugLogLevel = LogLevel(1)
	TraceLogLevel = LogLevel(2)
)

var logLevel = InfoLogLevel

func SetLogLevel(level LogLevel) {
	logLevel = level
}

func InfoLog(format string, v ...any) {
	if logLevel >= InfoLogLevel {
		log.Printf(format, v...)
	}
}

func DebugLog(format string, v ...any) {
	if logLevel >= DebugLogLevel {
		log.Printf("[debug]"+format, v...)
	}
}

func TraceLog(format string, v ...any) {
	if logLevel >= TraceLogLevel {
		log.Printf("[trace]"+format, v...)
	}
}
