package main

import (
	"log"

	flags "github.com/jessevdk/go-flags"
	voicepeakagent "github.com/yazawa-ichio/voicepeak-agent"
)

func main() {
	o := &opts{Port: 21952}
	_, err := flags.Parse(o)
	if err != nil {
		log.Fatal(err)
	}
	app := voicepeakagent.NewApp(o.AppPath, o.WorkDir)
	b := newBot(app)
	if err := b.start(o.Port); err != nil {
		log.Printf("bot error %v", err)
	}
}
