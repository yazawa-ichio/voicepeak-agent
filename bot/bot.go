package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	voicepeakagent "github.com/yazawa-ichio/voicepeak-agent"
)

type Bot struct {
	app *voicepeakagent.App
}

func NewBot(app *voicepeakagent.App) *Bot {
	return &Bot{app}
}

func (b *Bot) Start(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%+v", port), b)
}

func (b *Bot) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/get" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found " + r.URL.Path))
		return
	}
	contentType := r.Header.Get("Content-Type")
	var err error
	req := &voicepeakagent.SayRequest{}
	if contentType == "application/json" {
		d := json.NewDecoder(r.Body)
		err = d.Decode(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		query := r.URL.Query()
		if query.Has("text") {
			m := make(map[string]string)
			for key, _ := range query {
				m[key] = query.Get(key)
			}
			log.Print(m)
			req = &voicepeakagent.SayRequest{}
			buf, err := json.Marshal(m)
			log.Print(string(buf))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
			}
			err = json.Unmarshal(buf, &req)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			log.Print(req)
		} else {
			panic(r.URL.Path)
		}
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	req.FileName = fmt.Sprintf("tmp-%v.wav", time.Now().UnixMilli())
	path, err := b.app.Say(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer os.Remove(path)
	b.Play(path)
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (b *Bot) Play(path string) {
	log.Print(path)
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}
