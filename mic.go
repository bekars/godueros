package godueros

import (
	"os"
	"fmt"
	"time"
	"text/template"
	"github.com/gordonklaus/portaudio"
)

type DuMic struct {
	*portaudio.Stream
	buffer []float32
	index int
	duration int
}

func NewDuMic(delay time.Duration) (mic *DuMic) {
	h, err := portaudio.DefaultHostApi()
	CheckErr(err)
	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, nil/*h.DefaultOutputDevice*/)
	p.Input.Channels = 2

	mic = &DuMic{
		buffer: make([]float32, int(2 * p.SampleRate * delay.Seconds())),
		duration: int(delay.Seconds()),
		}
	mic.Stream, err = portaudio.OpenStream(p, mic.processAudio)
	CheckErr(err)
	return
}

func (mic *DuMic) processAudio(in, out []float32) {
	//for n := range in {
		//fmt.Printf("n: %d in: %f\n", n, in[n])
	//}

	for i := range in {
		mic.buffer[mic.index] = in[i]
		mic.index = (mic.index + 1) % len(mic.buffer)
	}
}

func (mic *DuMic)GetHw() {
	hs, err := portaudio.HostApis()
	CheckErr(err)
	err = tmpl.Execute(os.Stdout, hs)
	CheckErr(err)
}

var tmpl = template.Must(template.New("").Parse(
	`{{. | len}} host APIs: {{range .}}
	Name:                   {{.Name}}
	{{if .DefaultInputDevice}}Default input device:   {{.DefaultInputDevice.Name}}{{end}}
	{{if .DefaultOutputDevice}}Default output device:  {{.DefaultOutputDevice.Name}}{{end}}
	Devices: {{range .Devices}}
		Name:                      {{.Name}}
		MaxInputChannels:          {{.MaxInputChannels}}
		MaxOutputChannels:         {{.MaxOutputChannels}}
		DefaultLowInputLatency:    {{.DefaultLowInputLatency}}
		DefaultLowOutputLatency:   {{.DefaultLowOutputLatency}}
		DefaultHighInputLatency:   {{.DefaultHighInputLatency}}
		DefaultHighOutputLatency:  {{.DefaultHighOutputLatency}}
		DefaultSampleRate:         {{.DefaultSampleRate}}
	{{end}}
{{end}}`,
))

func (mic *DuMic)StartRecord() error {
	fmt.Printf("### Start Record %d Second\n", mic.duration)
	return mic.Start()
}

func (mic *DuMic)StopRecord() error {
	fmt.Printf("### Stop Record %d Second\n", mic.duration)
	return mic.Stop()
}

func (mic *DuMic)PlaySound() error {
	h, err := portaudio.DefaultHostApi()
	CheckErr(err)
	p := portaudio.LowLatencyParameters(nil, h.DefaultOutputDevice)
	p.Output.Channels = 2

	mic.Stream, err = portaudio.OpenStream(p, mic.playAudio)
	CheckErr(err)
	mic.index = 0
	return nil
}

func (mic *DuMic)playAudio(in, out []float32) {
	framePerBuffer := len(out)
	for i := 0; i < framePerBuffer; i++ {
		out[i] = mic.buffer[mic.index]
		mic.index = (mic.index + 1) % len(mic.buffer)
	}
}
