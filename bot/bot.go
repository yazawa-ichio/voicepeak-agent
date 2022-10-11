package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
	voicepeakagent "github.com/yazawa-ichio/voicepeak-agent"
	"github.com/yazawa-ichio/voicepeak-agent/speaker"
)

type bot struct {
	app     *voicepeakagent.App
	speaker *speaker.Speaker
}

func newBot(app *voicepeakagent.App) *bot {
	return &bot{app, speaker.NewSpeaker(app)}
}

func (b *bot) start(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/speak", b.speak)
	mux.HandleFunc("/help", b.help)
	voicepeakagent.InfoLog("start bot. http://localhost:%+v/speak?text=Hello", port)
	return http.ListenAndServe(fmt.Sprintf("localhost:%+v", port), mux)
}

func (b *bot) speak(w http.ResponseWriter, r *http.Request) {
	var err error
	d := schema.NewDecoder()
	req := &speakRequest{}
	err = d.Decode(req, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	say := req.ConvertSayRequest()
	voicepeakagent.InfoLog("say %+v", say)
	err = b.speaker.Play(say)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (b *bot) help(w http.ResponseWriter, r *http.Request) {
	help, err := b.app.GetHelp()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(help))
}
