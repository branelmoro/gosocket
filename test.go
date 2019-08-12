package main

import (
	"fmt"
	"gosocket"
	"io/ioutil"
	"os"
)

func main() {

	// 1) on malformed request
	gosocket.OnMalformedRequest = func(req gosocket.HttpRequest) {
		fmt.Println("OnMalformedRequest----\n", req.Raw())
		fmt.Println("OnMalformedRequest----\n", string(req.Raw()))
	}

	// 2) on http request
	gosocket.OnHttpRequest = func(w gosocket.HttpWriter, req gosocket.HttpRequest) {
		fmt.Println("OnHttpRequest----\n", string(req.Raw()))
		w.Close()
	}

	// 3) on websocket request
	gosocket.OnWebsocketRequest = func(w gosocket.HttpWriter, req gosocket.HttpRequest) {
		fmt.Println("OnWebsocketRequest----\n", string(req.Raw()) )
		err := w.UpgradeToWebsocket(nil)
		fmt.Println("UpgradeToWebsocket err----", err)
	}

	// 4) on websocket connection open
	gosocket.OnWebsocketOpen = func(w gosocket.WsWriter) {
		fmt.Println("OnWebsocketOpen--------------------")
		err := w.SendText("Hello, welcome to websocket protocol")
		fmt.Println("OnWebsocketOpen err is --------------------", err)
	}

	// 5) on message
	gosocket.OnMessage = func(w gosocket.WsWriter, msg gosocket.Message) {
		fmt.Println("OnMessage----", msg.Data())
	}

	// 5) on message
	gosocket.OnText = func(w gosocket.WsWriter, str string) {
		fmt.Println("OnText----Resending---", str)
		err := w.SendText(str)
		fmt.Println("Resending err is ----", err)
		// if str == "binary" {
		// 	w.SendBinary([]byte("message received:- " + str))
		// }
		// if str == "close" {
		// 	w.Close(nil)
		// }
	}

	// 5) on message
	gosocket.OnBinary = func(w gosocket.WsWriter, data []byte) {
		fmt.Println("OnBinary----", data)
	}

	// 6) on error
	gosocket.OnError = func(w gosocket.WsWriter, err error) {
		e, _ := err.(gosocket.Error)
		fmt.Println("OnError----", err, e.Error())
	}

	// 7) on connection close
	gosocket.OnClose = func(w gosocket.WsWriter, msg gosocket.CloseMsg) {
		fmt.Println("OnClose------", msg.Data(), msg.Code())
	}

	// 8) on ping
	gosocket.OnPing = func(w gosocket.WsWriter) {
		fmt.Println("OnPing----")
	}

	// 9) on pong
	gosocket.OnPong = func(w gosocket.WsWriter) {
		fmt.Println("OnPong----")
	}

	conf := gosocket.NewConf()

	content, err := ioutil.ReadFile("certs/server.pem")
	if err != nil {
		fmt.Println("Error----", err)
		os.Exit(1)
	}
	conf.CertPublic = content

	content, err = ioutil.ReadFile("certs/server.key")
	if err != nil {
		fmt.Println("Error----", err)
		os.Exit(1)
	}
	conf.CertPrivate = content

	fmt.Println("Conf is--------", conf);

	server := gosocket.NewServer(conf)

	server.Run();
}
