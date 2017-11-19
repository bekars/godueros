package godueros

import (
//	"fmt"
	"math"
	"github.com/gordonklaus/portaudio"
	"fmt"
)

const (
	Num_Seconds 		int32 = 5
	Sample_Rate 		int32 = 44100
	Frame_Per_Buffer 	int32 = 64
	Table_Size 			int = 126
)

type DuPlayer struct {
	sine []float32
	left_phase int
	right_phase int
	message string
	*portaudio.Stream
}

func NewDuPlayer() (player *DuPlayer, err error) {
	player = &DuPlayer{sine : make([]float32, Table_Size)}
	err = nil

	for i := 0; i < Table_Size; i++ {
		player.sine[i] = float32(math.Sin(float64(float32(i) / float32(Table_Size) * math.Pi * 2.)))
		fmt.Printf("%f\n", player.sine[i])
	}
	player.left_phase = 0
	player.right_phase = 0

	h, err := portaudio.DefaultHostApi()
	CheckErr(err)
	p := portaudio.LowLatencyParameters(nil, h.DefaultOutputDevice)
	p.Output.Channels = 2

	player.Stream, err = portaudio.OpenStream(p, player.processAudioCB)
	return
}

func (player *DuPlayer) processAudioCB(in, out []float32) {
	fmt.Printf("in/out len: %d/%d\n", len(in), len(out))
	for i := 0; i < len(out); i += 2 {
		out[i] = player.sine[player.left_phase]
		out[i+1] = player.sine[player.right_phase]
		player.left_phase += 3
		if (player.left_phase >= Table_Size) {
			player.left_phase -= Table_Size
		}

		player.right_phase += 9
		if (player.right_phase >= Table_Size) {
			player.right_phase -= Table_Size
		}
	}
}

