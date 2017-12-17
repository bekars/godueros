package godueros

import (
	"io"
	"math"
	"fmt"
	"os"
	"encoding/binary"
	"bytes"
	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/gordonklaus/portaudio"
	"time"
)

const (
	//Num_Seconds 		int32 = 5
	//Sample_Rate 		int32 = 44100
	//Frame_Per_Buffer 	int32 = 64
	Table_Size 			int = 126
)

type wavChunk struct {
	ChunkID 	  uint32  	// "RIFF"
	ChunkSize 	  uint32 	// Total bytes = 36 + Subchunk2Size
	Format 		  uint32   	// "WAVE"

	// sub-chunk "fmt"
	Subchunk1ID   uint32   	// "fmt "
	Subchunk1Size uint32 	// 16 for PCM
	AudioFormat   uint16   	// PCM = 1
	NumChannels   uint16  	// Mono = 1, Stereo = 2, etc.
	SampleRate    uint32    // 8000, 44100, etc.
	ByteRate      uint32  	// = SampleRate * NumChannels * BitsPerSample/8
	BlockAlign    uint16   	// = NumChannels * BitsPerSample/8
	BitsPerSample uint16 	// 8bits, 16bits, etc.

	// sub-chunk "data"
	Subchunk2ID   uint32 	// "data"
	Subchunk2Size uint32 	// data size
}

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

	chunk, err := readWavChunk(f)
	CheckErr(err)

	byte := uint32ToByte(chunk.ChunkID)
	fmt.Println(byte)
	fmt.Printf("%x %x %x %x\n", byte[0], byte[1], byte[2], byte[3])
	return err
}

type readerAtSeeker interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

func readWavChunk(r readerAtSeeker) (chunk *wavChunk, err error) {
	chunk = &wavChunk{}
	err = binary.Read(r, binary.BigEndian, chunk)
	if err != nil {
		return
	}
	return
}

func (speaker *DuSpeaker) PlayMP3File(filename string) (err error) {
	// create mpg123 decoder instance
	decoder, err := mpg123.NewDecoder("")
	CheckErr(err)

	CheckErr(decoder.Open(filename))
	defer decoder.Close()

	// get audio format information
	rate, channels, _ := decoder.GetFormat()

	// make sure output format does not change
	decoder.FormatNone()
	decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)

	out := make([]int16, 8192)
	stream, err := portaudio.OpenDefaultStream(0, channels, float64(rate), len(out), &out)
	CheckErr(err)
	defer stream.Close()

	CheckErr(stream.Start())
	defer stream.Stop()
	for {
		audio := make([]byte, 2 * len(out))
		n, err := decoder.Read(audio)
		fmt.Printf("### %d\n", n)
		if n > 0 {
			CheckErr(binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out))
			CheckErr(stream.Write())
		}
		if err == mpg123.EOF {
			time.Sleep(3 * time.Second) // wait voice play finish
			break
		}
		CheckErr(err)
	}
	return err
}

func (speaker *DuSpeaker) PlayMP3Audio(audio []byte) (err error) {
	// create mpg123 decoder instance
	decoder, err := mpg123.NewDecoder("")
	CheckErr(err)

	CheckErr(decoder.Feed(audio))

	// get audio format information
	rate, channels, _ := decoder.GetFormat()

	// make sure output format does not change
	decoder.FormatNone()
	decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)

	out := make([]int16, 8192)
	stream, err := portaudio.OpenDefaultStream(0, channels, float64(rate), len(out), &out)
	CheckErr(err)
	defer stream.Close()

	CheckErr(stream.Start())
	defer stream.Stop()
	for {
		audio := make([]byte, 2 * len(out))
		n, err := decoder.Read(audio)
		CheckErr(err)
		if n > 0 {
			CheckErr(binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out))
			CheckErr(stream.Write())
		}
		if err == mpg123.EOF {
			break
		}
	}
	return err
}
