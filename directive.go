package godueros

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"golang.org/x/net/http2"
//	"crypto/tls"
//	"net"
)

type DuDirective struct {

}

func (d *DuDirective)HTTP2() {
	req, _ := http.NewRequest("GET", "https://http2.golang.org/", nil)
	rt := &http2.Transport{
		/*
		AllowHTTP: true,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
		*/
	}

	client := &http.Client{
		Transport: rt,
	}
	req.Header.Add("Accept-Encoding", "identity")

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("http2 err: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	fmt.Println(res.Header)
}

func (d *DuDirective)ConnSrv() error {
	client := &http.Client{
		Transport: &http2.Transport{},
	}

	request, err := http.NewRequest("GET", "https://dueros-h2.baidu.com/dcs/v1/directives", nil)
	if err != nil {
		fmt.Println(err)
	}

	request.Header.Set("authorization", "Bearer " + access_token)
	request.Header.Set("dueros-device-id", device_id)

	fmt.Println("Do request")
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Get error", err)
	}

	go func() {
		fmt.Println("Read body")
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Get body error", err)
		}

		fmt.Println(string(body))
		fmt.Println(response.Header)
		response.Body.Close()
	}()

	return err
}

