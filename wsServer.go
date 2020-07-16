package main

import (
	"fmt"
	"gosocket"
	"io/ioutil"
	"os"
)

func main() {

	sConf := gosocket.NewServerConf()
	// 1) on malformed request
	sConf.OnMalformedRequest = func(req gosocket.HttpRequest) {
		fmt.Println("OnMalformedRequest----\n", req.Raw())
		fmt.Println("OnMalformedRequest----\n", string(req.Raw()))
	}

	// 2) on http request
	sConf.OnHttpRequest = func(w gosocket.HttpWriter, req gosocket.HttpRequest) {
		fmt.Println("OnHttpRequest----\n", string(req.Raw()))
		w.Close()
	}

	// 3) on websocket request
	sConf.OnWebsocketRequest = func(w gosocket.HttpWriter, req gosocket.HttpRequest) {
		fmt.Println("OnWebsocketRequest----\n", string(req.Raw()) )
		err := w.UpgradeToWebsocket(nil)
		fmt.Println("UpgradeToWebsocket err----", err)
	}


	wsConf := NewWsConf()

	// 4) on websocket connection open
	wsConf.OnWebsocketOpen = func(w gosocket.WsWriter) {
		fmt.Println("OnWebsocketOpen--------------------")
		err := w.SendText("Hello, welcome to websocket protocol")
		fmt.Println("OnWebsocketOpen err is --------------------", err)
	}

	// 5) on message
	wsConf.OnMessage = func(w gosocket.WsWriter, msg gosocket.Message) {
		fmt.Println("OnMessage----", msg.Data())
	}

	// 5) on message
	wsConf.OnText = func(w gosocket.WsWriter, str string) {
		fmt.Println("OnText----Resending---", str)
		err := w.SendText(str)
		fmt.Println("Resending err is ----", err)
	}

	// 5) on message
	wsConf.OnBinary = func(w gosocket.WsWriter, data []byte) {
		fmt.Println("OnBinary----", data)
	}

	// 6) on error
	wsConf.OnError = func(w gosocket.WsWriter, err error) {
		e, _ := err.(gosocket.Error)
		fmt.Println("OnError----", err, e.Error())
	}

	// 7) on connection close
	wsConf.OnClose = func(w gosocket.WsWriter, msg gosocket.CloseMsg) {
		fmt.Println("OnClose------", msg.Data(), msg.Code())
	}

	// 8) on ping
	wsConf.OnPing = func(w gosocket.WsWriter) {
		fmt.Println("OnPing----")
	}

	// 9) on pong
	wsConf.OnPong = func(w gosocket.WsWriter) {
		fmt.Println("OnPong----")
	}

	// add wsConf to server conf
	sConf.WsConf = wsConf

	content, err := ioutil.ReadFile("certs/server.pem")
	if err != nil {
		fmt.Println("Error----", err)
		os.Exit(1)
	}
	sConf.CertPublic = content

	content, err = ioutil.ReadFile("certs/server.key")
	if err != nil {
		fmt.Println("Error----", err)
		os.Exit(1)
	}
	sConf.CertPrivate = content

	fmt.Println("sConf is--------", sConf);

	server := gosocket.NewServer(sConf)

	server.Run();
}
