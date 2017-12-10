package godueros

import (
	"io"
	"math"
	"github.com/gordonklaus/portaudio"
	"fmt"
	"os"
	"encoding/binary"
)

const (
	//Num_Seconds 		int32 = 5
	//Sample_Rate 		int32 = 44100
	//Frame_Per_Buffer 	int32 = 64
	Table_Size 			int = 126
)

type DuSpeaker struct {
	sine []float32
	left_phase int
	right_phase int
	message string
	*portaudio.Stream
}

func NewDuSpeaker() (player *DuSpeaker, err error) {
	player = &DuSpeaker{sine : make([]float32, Table_Size)}
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

func (player *DuSpeaker) processAudioCB(in, out []float32) {
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

func (player *DuSpeaker) PlayFile(filename string) error {
	f, err := os.Open(filename)
	CheckErr(err)
	defer f.Close()

	id, data, err := readChunk(f)
	CheckErr(err)
	if id.String() != "FORM" {
		fmt.Println("bad file format")
		return err
	}
	_, err = data.Read(id[:])
	CheckErr(err)
	if id.String() != "AIFF" {
		fmt.Println("bad file format")
		return err
	}

	return err
}

type readerAtSeeker interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

type commonChunk struct {
	NumChans      int16
	NumSamples    int32
	BitsPerSample int16
	SampleRate    [10]byte
}
type ID [4]byte

func (id ID) String() string {
	return string(id[:])
}

func readChunk(r readerAtSeeker) (id ID, data *io.SectionReader, err error) {
	_, err = r.Read(id[:])
	if err != nil {
		return
	}
	var n int32
	err = binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return
	}
	fmt.Printf("## ReadChunk %d\n", n)
	off, _ := r.Seek(0, 1)
	data = io.NewSectionReader(r, off, int64(n))
	_, err = r.Seek(int64(n), 1)
	return
}
