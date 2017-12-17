package main

import (
	"time"
	"github.com/bekars/godueros"
	"github.com/gordonklaus/portaudio"
	"fmt"
)

const (
	NUM_SEC = 1
)

func playFile() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	mic := godueros.NewDuMic(time.Second * NUM_SEC)
	mic.LoadFile("audio.wav")
	mic.PlaySound()
	mic.StartRecord()
	time.Sleep(time.Second * NUM_SEC)
	mic.StopRecord()
}

func recordFile() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	fmt.Println("Start Record Voice")
	mic := godueros.NewDuMic(time.Second * NUM_SEC)
	mic.GetHw()
	mic.StartRecord()
	time.Sleep(time.Second * NUM_SEC)
	mic.StopRecord()

	mic.WriteFile()
	mic.PlaySound()
	mic.StartRecord()
	time.Sleep(time.Second * NUM_SEC)
	mic.StopRecord()
}

func main() {
	//playFile()
	//return

	//recordFile()
	//return

	//directive := &godueros.DuDirective{}
	//directive.ConnSrv()

	//event := &godueros.DuEvent{}
	//event.SendEvent()
	//time.Sleep(time.Second * 5)

	core, _ := godueros.NewDuCore()
	core.Run()
	return
}

