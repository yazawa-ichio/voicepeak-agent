package voicepeakagent

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type App struct {
	path, workpace string
	mutex          *sync.Mutex
}

func NewApp(path, workpace string) *App {
	return &App{
		path:     path,
		workpace: workpace,
		mutex:    &sync.Mutex{},
	}
}

func (p *App) run(args ...string) string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	defer newStopwatch("run").Stop()
	TraceLog("run %s", args)
	ret, err := exec.Command(p.path, args...).Output()
	if err != nil {
		log.Fatalf("run fail %s", err)
	}
	return string(ret)
}

func (p *App) tryRun(args ...string) (string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	defer newStopwatch("run").Stop()
	TraceLog("run %s", args)
	ret, err := exec.Command(p.path, args...).Output()
	return string(ret), err
}

func (p *App) Help() {
	defer newStopwatch("help").Stop()
	ret := p.run("--help")
	InfoLog("%s", ret)
}

func (p *App) GetHelp() (string, error) {
	defer newStopwatch("help").Stop()
	ret, err := p.tryRun("--help")
	TraceLog("%s", ret)
	return ret, err
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

func (p *App) SimpleSay(fileName, text string, narrator Narrator) string {
	defer newStopwatch("say").Stop()
	path := filepath.Join(p.workpace, fileName)
	p.run("-s", text, "-n", narrator.String(), "-o", path)
	DebugLog("output %s", path)
	return path
}

func (p *App) SayToFile(path string, req *SayRequest) (string, error) {
	defer newStopwatch("say to file").Stop()
	args, err := req.GetArgs(path)
	if err != nil {
		return "", err
	}
	_, err = p.tryRun(args...)
	if err != nil {
		return "", err
	}
	DebugLog("output %s", path)
	return path, nil
}

func (p *App) SayToFileByJson(path string, payload []byte) (string, error) {
	defer newStopwatch("say to file").Stop()
	req := &SayRequest{}
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return "", err
	}
	args, err := req.GetArgs(path)
	if err != nil {
		return "", err
	}
	_, err = p.tryRun(args...)
	if err != nil {
		return "", err
	}
	DebugLog("output %s", path)
	return path, nil
}

func (p *App) Say(req *SayRequest, f func(path string) error) error {
	defer newStopwatch("say").Stop()
	path := filepath.Join(p.workpace, fmt.Sprintf("tmp-%v.wav", time.Now().UnixMilli()))
	args, err := req.GetArgs(path)
	if err != nil {
		return err
	}
	_, err = p.tryRun(args...)
	if err != nil {
		return err
	}
	DebugLog("output %s", path)
	defer os.Remove(path)
	err = f(path)
	if err != nil {
		return err
	}
	return nil
}

func (p *App) SayByJson(payload []byte, f func(path string) error) error {
	defer newStopwatch("say").Stop()
	req := &SayRequest{}
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return err
	}
	path := filepath.Join(p.workpace, fmt.Sprintf("tmp-%v.wav", time.Now().UnixMilli()))
	args, err := req.GetArgs(path)
	if err != nil {
		return err
	}
	_, err = p.tryRun(args...)
	if err != nil {
		return err
	}
	DebugLog("output %s", path)
	defer os.Remove(path)
	err = f(path)
	if err != nil {
		return err
	}
	return nil
}
