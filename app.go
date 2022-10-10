package voicepeakagent

import (
	"encoding/json"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

type App struct {
	path, workpace string
}

func NewApp(path, workpace string) *App {
	return &App{
		path:     path,
		workpace: workpace,
	}
}

func (p *App) run(args ...string) string {
	TraceLog("run %s", args)
	ret, err := exec.Command(p.path, args...).Output()
	if err != nil {
		log.Fatalf("run fail %s", err)
	}
	return string(ret)
}

func (p *App) Help() {
	defer newStopwatch("help").Stop()
	ret := p.run("--help")
	InfoLog("%s", ret)
}

func (p *App) ListNarrator() []Narrator {
	defer newStopwatch("list-narrator").Stop()
	ret := p.run("--list-narrator")
	var list []Narrator
	for _, n := range strings.Split(strings.ReplaceAll(ret, "\r\n", "\n"), "\n") {
		if len(strings.TrimSpace(n)) > 0 {
			list = append(list, Narrator(n))
		}
	}
	DebugLog("%s", list)
	return list
}

func (p *App) ListEmotion(narrator Narrator) []string {
	defer newStopwatch("list-emotion").Stop()
	ret := p.run("--list-emotion", string(narrator))
	list := strings.Split(strings.ReplaceAll(ret, "\r\n", "\n"), "\n")
	DebugLog("%s %s", narrator, list)
	return list
}

func (p *App) SimpleSay(text string, narrator Narrator) string {
	defer newStopwatch("say").Stop()
	name := text + "-" + narrator.String() + ".wav"
	p.run("-s", text, "-n", narrator.String(), "-o", filepath.Join(p.workpace, name))
	DebugLog("output %s", name)
	return filepath.Join(p.workpace, name)
}

func (p *App) Say(req *SayRequest) (string, error) {
	defer newStopwatch("say").Stop()
	path := filepath.Join(p.workpace, req.FileName)
	args, err := req.GetArgs(p.workpace)
	if err != nil {
		return "", err
	}
	p.run(args...)
	DebugLog("output %s", req.FileName)
	return path, nil
}

func (p *App) SayByJson(payload []byte) (string, error) {
	defer newStopwatch("say").Stop()
	req := &SayRequest{}
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return "", err
	}
	path := filepath.Join(p.workpace, req.FileName)
	args, err := req.GetArgs(p.workpace)
	if err != nil {
		return "", err
	}
	p.run(args...)
	DebugLog("output %s", req.FileName)
	return path, nil
}
