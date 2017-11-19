package godueros

import (
	"github.com/gordonklaus/portaudio"
	"time"
//	"fmt"
	"os"
	"text/template"
)

type DuMic struct {
	hw int
	*portaudio.Stream
	buffer []float32
	i int
}

func NewDuMic(delay time.Duration) (mic *DuMic) {
	h, err := portaudio.DefaultHostApi()
	CheckErr(err)
	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, h.DefaultOutputDevice)
	p.Input.Channels = 1
	p.Output.Channels = 1

	mic = &DuMic{buffer: make([]float32, int(p.SampleRate * delay.Seconds()))}
	mic.Stream, err = portaudio.OpenStream(p, mic.processAudio)
	CheckErr(err)
	return
}

func (mic *DuMic) processAudio(in, out []float32) {
	//for n := range in {
		//fmt.Printf("n: %d in: %f\n", n, in[n])
	//}

	for i := range in {
		//out[i] = .7 * mic.buffer[mic.i]
		mic.buffer[mic.i] = in[i]
		mic.i = (mic.i + 1) % len(mic.buffer)
	}
}

func (mic *DuMic)SetHw(hw int) {
	mic.hw = hw
}

func (mic *DuMic)GetHw() int {
	hs, err := portaudio.HostApis()
	CheckErr(err)
	err = tmpl.Execute(os.Stdout, hs)
	CheckErr(err)
	return mic.hw
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
	return nil
}

func (mic *DuMic)PlaySound() error {
	h, err := portaudio.DefaultHostApi()
	CheckErr(err)
	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, h.DefaultOutputDevice)
	p.Input.Channels = 1
	p.Output.Channels = 1

	mic.Stream, err = portaudio.OpenStream(p, mic.playAudio)
	CheckErr(err)
	return nil
}

func (mic *DuMic) playAudio(in, out []float32) {

	for i := range mic.buffer {
		out[i] = mic.buffer[i]
	}
}
