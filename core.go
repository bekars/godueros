package godueros

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"net/http"
)

type DuStage int
const (
	DU_WAKEUP  DuStage = iota
	DU_EVENT
	DU_PLAY
)

type DuCore struct {
	Client     *http.Client
	Stage       DuStage
	recorderCh  chan VoiceChannel
	directive  *DuDirective
	event      *DuEvent
	recorder   *DuRecorder
	Speaker    *DuSpeaker
}

func NewDuCore()(core *DuCore, err error)  {
	err = nil
	core = &DuCore{
		directive: 	&DuDirective{},
		event: 		&DuEvent{},
		Speaker: 	&DuSpeaker{},
	}

	return core, err
}

func (core *DuCore)Run() (err error) {
	// init portaudio
	portaudio.Initialize()
	defer portaudio.Terminate()

	// wake up me
	core.Stage = DU_WAKEUP

	core.Stage = DU_EVENT
	// new http2 client
	core.Client, err = core.directive.ConnSrv()

	// new event
	core.event, _ = NewDuEvent(core)

	// new recorder
	core.recorderCh = make(chan VoiceChannel, 128)
	core.recorder, err = NewDuRecorder(core.recorderCh)
	core.recorder.Start()

	fmt.Println("Please say something ...")

	select {
	case <- core.recorderCh:
		b := make([]byte, 16000)
		for n := range core.recorderCh {
			for i := range n.Buffer {
				//fmt.Printf("%d,", n.Buffer[i])

				bytes := int16ToByte(n.Buffer[i])
				for m := range bytes {
					b = append(b, bytes[m])
				}
			}
			if core.recorder.Rs == RS_TIMEOUT {
				//fmt.Println("Receive Voice Timeout")
				break
			}
		}
		core.recorder.Stop()

		core.event.SendDCS(b)

		core.Speaker.PlayMP3File(REPLY_FILE)
	}

	return err
}

func (core *DuCore) GetRecorderCh() (ch chan VoiceChannel) {
	return core.recorderCh
}
