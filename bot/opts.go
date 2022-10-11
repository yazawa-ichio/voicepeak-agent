package main

type opts struct {
	Port    int    `short:"p" long:"port" description:"local server port"`
	AppPath string `short:"a" long:"app" description:"VOICEPEAK app path" required:"true"`
	WorkDir string `short:"d" long:"dir" description:"voice workspace" required:"true"`
}
