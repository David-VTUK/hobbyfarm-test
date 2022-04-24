package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func main() {

	pass, err := testDNSResolution()
	if err != nil {
		fmt.Println("Unable to resolve sslip.io address")
	}

	if pass == true {
		fmt.Println("DNS resolution pass")
	}

	pass, err = testSslipConnection()
	if err != nil {
		fmt.Println("Unable to connect to sslip.io", err)
	}

	if pass == true {
		fmt.Println("HTTP request to sslip.io pass")
	}

	pass, err = testWebsocketConnection()
	if err != nil {
		fmt.Println("Unable to initiate websocket connection", err)
	}

	if pass == true {
		fmt.Println("Websocket connection pass")
	}
}

func testDNSResolution() (bool, error) {
	ips, err := net.LookupIP("104.155.144.4.sslip.io")

	if err != nil {
		return false, err
	}

	if ips[0].String() != "104.155.144.4" {
		return false, errors.New("unexpected result resolving sslip.io address")
	}

	return true, nil
}

func testSslipConnection() (bool, error) {
	resp, err := http.Get("https://sslip.io")
	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	found := strings.Contains(string(body), "Welcome to sslip.io")

	if found == true {
		return true, nil
	} else {
		return false, errors.New("unexpected result accessing https://sslip.io")
	}
}

func testWebsocketConnection() (bool, error) {

	c, _, err := websocket.DefaultDialer.Dial("wss://ws.ifelse.io", nil)

	defer c.Close()

	if err != nil {
		return false, err
	}

	_, response, err := c.ReadMessage()
	if err != nil || response == nil {
		return false, err
	}

	found := strings.Contains(string(response), "Request served by")

	if found == true {
		return true, nil
	}

	return true, nil
}
