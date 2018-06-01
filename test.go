package main

import (
	"fmt"
	"net"
	"gosocket"
)

func main() {

	// 1) on connection open
	gosocket.OnWebsocketOpen = func(conn net.Conn) {

		fmt.Println("in OnWebsocketOpen------------")

	}


	// 2) on message
	gosocket.OnMessage(conn gosocket.Conn, a []byte) {

		fmt.Println("in OnMessage------------")
	}

	// 3) on error
	gosocket.OnError(conn net.Conn) {

		fmt.Println("in OnMessage------------")
	}

	// 4) on connection close
	gosocket.OnClose(conn gosocket.Conn) {

		fmt.Println("in OnClose------------")

	}


	gosocket.StartServer(config)
	
}
