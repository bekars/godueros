package main

import (
/*	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/http2"
*/
	"github.com/bekars/godueros"
	"fmt"
	"time"
	"github.com/gordonklaus/portaudio"
)

func main() {

	portaudio.Initialize()
	defer portaudio.Terminate()

	mic := godueros.NewDuMic(time.Second / 3)
	fmt.Printf("hw %d\n", mic.GetHw())
	defer mic.Close()
	CheckErr(mic.Start())
	time.Sleep(60 * time.Second)
	CheckErr(mic.Stop())

	/*
	godueros.Abc()

	fmt.Printf("%s Hello!", "DuerOS")

	client := http.Client{
		Transport: &http2.Transport{},
	}

	request, err := http.NewRequest("GET", "https://dueros-h2.baidu.com/dcs/v1/directives", nil)
	if err != nil {
		fmt.Println(err)
	}

	request.Header.Set("authorization", "Bearer 23.34e01113d7efe0c1ea0ad11a03aee31e.2592000.1512300410.2300547068-10218111")
	request.Header.Set("dueros-device-id", "wh9foSD4KnZAa8kyG1l62SdwkNHlVfuH")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Get error", err)
	}

	ping_request, err := http.NewRequest("GET", "https://dueros-h2.baidu.com/dcs/v1/ping", nil)
	if err != nil {
		fmt.Println(err)
	}
	ping_request.Header.Set("authorization", "Bearer 23.34e01113d7efe0c1ea0ad11a03aee31e.2592000.1512300410.2300547068-10218111")
	ping_request.Header.Set("dueros-device-id", "wh9foSD4KnZAa8kyG1l62SdwkNHlVfuH")

	ping_response, err := client.Do(ping_request)

	ping_body, err := ioutil.ReadAll(ping_response.Body)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Get body error", err)
	}

	fmt.Println(string(ping_body))
	fmt.Println(ping_response.Header)


	event_request, err := http.NewRequest("POST", "https://dueros-h2.baidu.com/dcs/v1/events", nil)
	event_request.Header.Set("authorization", "Bearer 23.34e01113d7efe0c1ea0ad11a03aee31e.2592000.1512300410.2300547068-10218111")
	event_request.Header.Set("dueros-device-id", "wh9foSD4KnZAa8kyG1l62SdwkNHlVfuH")
	event_request.Header.Set("content-type", "multipart/form-data; boundary=this-is-a-boundary")
	event_response, err := client.Do(event_request)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Get error", err)
	}

	event_body, err := ioutil.ReadAll(event_response.Body)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Get body error", err)
	}

	fmt.Println(string(event_body))
	fmt.Println(event_response.Header)

	defer response.Body.Close()
	fmt.Println(response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Get body error", err)
	}

	fmt.Println(string(body))
*/

}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}