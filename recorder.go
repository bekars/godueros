/**
 * "AUDIO_L16_RATE_16000_CHANNELS_1": 16bit线性PCM音频，16kHz采样率，单声道，Little endian byte order
 */
package godueros

import (
	"github.com/gordonklaus/portaudio"
)

type RecorderStatus int

const (
	RS_RECORDING RecorderStatus = iota
	RS_TIMEOUT
	RS_WAITING
)

type VoiceChannel struct {
	Buffer []int16
}

type DuRecorder struct {
	*portaudio.Stream
	buffer []int16
	len int
	ch chan<- VoiceChannel
	chunk_num int
	Rs RecorderStatus
}

func NewDuRecorder(ch chan<- VoiceChannel) (recorder *DuRecorder, err error) {
	h, err := portaudio.DefaultHostApi()
	CheckErr(err)

	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, nil)
	p.Input.Channels = int(Channels_Num)
	p.SampleRate = float64(Sample_Rate)
	p.FramesPerBuffer = int(Sample_Rate / Chunk_Per_Sec)

	recorder = &DuRecorder{
		buffer: make([]int16, p.Input.Channels * p.FramesPerBuffer),
		len: p.Input.Channels * p.FramesPerBuffer,
		chunk_num: 0,
		Rs: RS_RECORDING,
	}
	recorder.ch = ch
	recorder.Stream, err = portaudio.OpenStream(p, recorder.recordVoice)
	CheckErr(err)
	return
}

func (recorder *DuRecorder) recordVoice(in, out []int16) {
	for i := range in {
		if i >= recorder.len {
			break
		}
		recorder.buffer[i] = in[i]
	}

	vc := VoiceChannel{
		Buffer: make([]int16, len(in)),
	}
	for i := range in {
		vc.Buffer[i] = in[i]
	}

	recorder.chunk_num++
	if recorder.chunk_num >= Timeout_Chunk_Num {
		recorder.Rs = RS_TIMEOUT
	}
	recorder.ch <- vc
}

