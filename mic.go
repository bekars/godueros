package godueros

import (
	"github.com/gordonklaus/portaudio"
	"time"
	"fmt"
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
	p.Input.Channels = 2
	p.Output.Channels = 2

	mic = &DuMic{buffer: make([]float32, int(p.SampleRate*delay.Seconds()))}
	mic.Stream, err = portaudio.OpenStream(p, mic.processAudio)
	CheckErr(err)
	return
}

func (mic *DuMic) processAudio(in, out []float32) {
	for n := range in {
		fmt.Printf("n: %d in: %f\n", n, in[n])
	}

	for i := range out {
		out[i] = .7 * mic.buffer[mic.i]
		mic.buffer[mic.i] = in[i]
		mic.i = (mic.i + 1) % len(mic.buffer)
	}
}

func (mic *DuMic)SetHw(hw int) {
	mic.hw = hw
}

func (mic *DuMic)GetHw() int {
	return mic.hw
}

func (mic *DuMic)StartRecord() error {
	return nil
}