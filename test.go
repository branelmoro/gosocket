package main

import (
	"fmt"
	"gosocket"
)

func main() {

	// 1) on connection open
	gosocket.OnWebsocketOpen = func(conn *gosocket.Conn, a *[]byte) bool {
		data := *a
		fmt.Println("Opening Websocket")
		fmt.Println(string(data))
		// conn.Write([]byte("Request data received ------ "))
		// conn.Write(data)
		return true
	}


	// 2) on message
	gosocket.OnMessage = func(conn *gosocket.Conn, a *[]byte) {
		data := *a
		fmt.Println(string(data))
		conn.Write([]byte("Message received ------ "))
		conn.Write(data)
	}

	// 3) on error
	gosocket.OnError = func(conn *gosocket.Conn) {

		fmt.Println("in OnMessage------------")
	}

	// 4) on connection close
	gosocket.OnClose = func(conn *gosocket.Conn) {

		fmt.Println("in OnClose------------")

	}


	gosocket.StartServer()
	
}
