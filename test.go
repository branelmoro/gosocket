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
		fmt.Println(string(data), data)
		// conn.Write([]byte("Request data received ------ "))
		// conn.Write(data)
		return true
	}

	// 2) on message
	gosocket.OnMessage = func(conn *gosocket.Conn, a *gosocket.Message) {
		// message := *a
		msg := a.GetData()
		fmt.Println("in OnMessage - ", msg)
		data := *msg
		fmt.Println(string(data),data)

		r_data := []byte("Message received ------ " + string(data))
		conn.WriteMessage(&r_data)
		// conn.Write(data)
	}

	// 3) on error
	gosocket.OnError = func(conn *gosocket.Conn) {

		fmt.Println("in OnMessage------------")
	}

	// 4) on connection close
	gosocket.OnClose = func(conn *gosocket.Conn, a *gosocket.Message) {

		fmt.Println("in OnClose------------")

	}


	gosocket.StartServer()
	
}
