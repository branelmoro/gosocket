package main

import (
	"fmt"
	"gosocket"
)

func main() {

	// 1) on connection open
	gosocket.OnWebsocketOpen = func(conn gosocket.Conn) {

		fmt.Println("in OnWebsocketOpen------------")

	}


	// 2) on message
	gosocket.OnMessage = func(conn gosocket.Conn, a []byte) {

		fmt.Println("in OnMessage------------", string(a))

		conn.Write([]byte("You sent ------ "))
		conn.Write(a)


	}

	// 3) on error
	gosocket.OnError = func(conn gosocket.Conn) {

		fmt.Println("in OnMessage------------")
	}

	// 4) on connection close
	gosocket.OnClose = func(conn gosocket.Conn) {

		fmt.Println("in OnClose------------")

	}


	gosocket.StartServer()
	
}
