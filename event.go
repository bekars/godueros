package godueros

import (
	"fmt"
	"io/ioutil"
//	"log"
	"net/http"
	"golang.org/x/net/http2"
	"os"
	"bytes"
)

var (
	boundary = "--this-is-a-boundary"
	json = `{
	"clientContext": [
	    0,
	    0,
	    30,
	    0
	],
	"event": {
	    "header": {
	        "namespace": "ai.dueros.device_interface.voice_input",
	        "name": "ListenStarted",
	        "messageId": "messageId-111111",
	        "dialogRequestId": "123456"
	    },
	    "payload": {
	        "format": "AUDIO_L16_RATE_16000_CHANNELS_1"
	    }
	}
}`
)

type DuEvent struct {
	client *http.Client
	core *DuCore
	ch chan VoiceChannel
}
func NewDuEvent(core *DuCore) (event *DuEvent, err error) {
	return newDuEvent(core.Client, core.GetRecorderCh(), core)
}

func newDuEvent(client *http.Client, ch chan VoiceChannel, core *DuCore) (event *DuEvent, err error) {
	event = &DuEvent{
		client: client,
		ch: ch,
		core: core,
	}
	return event, err
}

func (e *DuEvent) SendDCS(voice []byte) (err error) {
	client := e.client
	dcs := &DuDCS{}
	multipart, _ := dcs.GetMultiPartData(json, voice, boundary)

	request, err := http.NewRequest("POST", "https://dueros-h2.baidu.com/dcs/v1/events", multipart)
	if err != nil {
		fmt.Println(err)
	}

	request.Header.Set("authorization", "Bearer " + access_token)
	request.Header.Set("dueros-device-id", device_id)
	request.Header.Set("content-type", "multipart/form-data; boundary=" + boundary)


	fmt.Println("Do Event Request")
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Get Event Error: ", err)
	}
	defer response.Body.Close()

	fmt.Println("Read Event Body")
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Get Event Body Error: ", err)
	}

	fmt.Println(response.Header)
	dcs.ReadMultiPartData(bytes.NewReader(body), "___dueros_dcs_v1_boundary___")

	return err
}

func (e *DuEvent) getVoiceFromFile(filename string) []byte {
	in, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Open audio file err: %v\n", err)
	}
	defer in.Close()

	bytes, err := ioutil.ReadAll(in)
	if err != nil {
		fmt.Printf("Read audio file err: %v\n", err)
	}
	return bytes
}

func (e *DuEvent) SendEvent() {
	client := &http.Client{
		Transport: &http2.Transport{},
	}

	dcs := &DuDCS{}
	multipart, _ := dcs.GetMultiPartData(json, e.getVoiceFromFile("voice/stock.wav"), boundary)

	fmt.Println(multipart.String())

	request, err := http.NewRequest("POST", "https://dueros-h2.baidu.com/dcs/v1/events", multipart)
	if err != nil {
		fmt.Println(err)
	}

	request.Header.Set("authorization", "Bearer " + access_token)
	request.Header.Set("dueros-device-id", device_id)
	request.Header.Set("content-type", "multipart/form-data; boundary=" + boundary)

	fmt.Println("Do Event Request")
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Get Event Error: ", err)
	}

	fmt.Println("Read Event Body")
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Get Event Body Error: ", err)
	}

	//fmt.Println(string(body))
	fmt.Println(response.Header)
	dcs.ReadMultiPartData(bytes.NewReader(body), "___dueros_dcs_v1_boundary___")
	defer response.Body.Close()
}

