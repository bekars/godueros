/**
 * "AUDIO_L16_RATE_16000_CHANNELS_1": 16bit线性PCM音频，16kHz采样率，单声道，Little endian byte order
 */
package godueros

import (
	"os"
	"fmt"
	"time"
	"text/template"
	"github.com/gordonklaus/portaudio"
	"bufio"
	"io/ioutil"
)

type DuMic struct {
	*portaudio.Stream
	buffer []int16
	//buffer []float32
	index int
	duration int
}

func NewDuMic(delay time.Duration) (mic *DuMic) {
	h, err := portaudio.DefaultHostApi()
	CheckErr(err)
	//p := portaudio.HighLatencyParameters(h.DefaultInputDevice,nil)
	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, nil/*h.DefaultOutputDevice*/)
	p.Input.Channels = int(Channels_Num)
	p.SampleRate = float64(Sample_Rate)
	p.FramesPerBuffer = Frame_Per_Buffer

	mic = &DuMic{
		//buffer: make([]float32, int(int(p.Input.Channels) * int(p.SampleRate) * int(delay.Seconds()))),
		buffer: make([]int16, int(int(p.Input.Channels) * int(p.SampleRate) * int(delay.Seconds()))),
		index: 0,
		duration: int(delay.Seconds()),
	}
	mic.Stream, err = portaudio.OpenStream(p, mic.processAudio)
	CheckErr(err)
	return
}

func (mic *DuMic) processAudio(in, out []int16) {
	fmt.Printf("### Audio in len: %d\n", len(in))
	for i := range in {
		mic.buffer[mic.index] = in[i]
		mic.index = (mic.index + 1) % len(mic.buffer)
	}
}

/*
func (mic *DuMic) processAudio(in, out []float32) {
	fmt.Printf("### Audio in len: %d\n", len(in))
	for i := range in {
		//mic.buffer[mic.index] = in[i]
		mic.buffer[mic.index] = int16(in[i])
		mic.index = (mic.index + 1) % len(mic.buffer)
	}
}
*/

func (mic *DuMic)WriteFile() {
	out, err := os.OpenFile("audio.wav", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("Open audio file err: %v\n", err)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	for i := 0; i < mic.index; i++ {
		//bytes := float32ToByte(mic.buffer[i])
		bytes := int16ToByte(mic.buffer[i])
		n, err := writer.Write(bytes)
		if n == 0  {
			CheckErr(err)
			break
		}
	}
	writer.Flush()
}

func (mic *DuMic)LoadFile(filename string) {
	in, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Open audio file err: %v\n", err)
	}
	defer in.Close()

	bytes, err := ioutil.ReadAll(in)
	if err != nil {
		fmt.Printf("Read audio file err: %v\n", err)
	}

	var index int
	for i:=0; i<len(bytes); i=i+2 {
		mic.buffer[index] = byteToInt16(bytes[i:])
		//mic.buffer[index] = byteToFloat32(bytes[i:])
		index += 1
		if index >= len(mic.buffer) {
			break
		}
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
	//p := portaudio.HighLatencyParameters(nil, h.DefaultOutputDevice)
	p := portaudio.LowLatencyParameters(nil, h.DefaultOutputDevice)
	p.Output.Channels = int(Channels_Num)
	p.SampleRate = float64(Sample_Rate)
	p.FramesPerBuffer = Frame_Per_Buffer

	mic.index = 0
	mic.Stream, err = portaudio.OpenStream(p, mic.playAudio)
	CheckErr(err)
	return nil
}

func (mic *DuMic)playAudio(in, out []int16) {
	fmt.Printf("### Audio out len: %d\n", len(out))
	for i := 0; i < len(out); i++ {
		out[i] = mic.buffer[mic.index]
		mic.index = (mic.index + 1) % len(mic.buffer)
	}
}


/*
func (mic *DuMic)playAudio(in, out []float32) {
	fmt.Printf("### Audio out len: %d\n", len(out))
	framePerBuffer := len(out)
	for i := 0; i < framePerBuffer; i++ {
		out[i] = float32(mic.buffer[mic.index])
		mic.index = (mic.index + 1) % len(mic.buffer)
	}
}
*/
