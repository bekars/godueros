package godueros

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestWriteMultiPartData(t *testing.T) {
	dcs := &DuDCS{}

	boundary := "--- ThisIsABoundary ---"
	json := `"clientContext": [
    {{ai.dueros.device_interface.alerts.AlertsState}},
    {{ai.dueros.device_interface.audio_player.PlaybackState}},
    {{ai.dueros.device_interface.speaker_controller.VolumeState}},
    {{ai.dueros.device_interface.voice_output.SpeechState}}
],
"event": {
    "header": {
        "namespace": "ai.dueros.device_interface.voice_input",
        "name": "ListenStarted",
        "messageId": "{{STRING}}",
        "dialogRequestId": "{{STRING}}"
    },
    "payload": {
        "format": "{{STRING}}"
    }
}`
	audio := []byte("AUDIO_SOUND")

	multipart, _ := dcs.GetMultiPartData(json, audio, boundary)
	fmt.Printf("### multipart: \n%s\n", multipart.String())
	assert.NotEqual(t, "", multipart, "multipart shouldn't empty")
}
