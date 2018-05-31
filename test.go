import "gosocket"



func main() {

	// 1) on connection open

	gosocket.OnWebsocketOpen(func(conn gosocket.Conn) {


	})


	// 2) on message

	gosocket.OnMessage(conn gosocket.Conn, a byte[]) {


	})

	// 3) on error

	gosocket.OnError(conn gosocket.Conn) {


	})

	// 4) on connection close


	gosocket.OnClose(conn gosocket.Conn) {


	})


	gosocket.StartServer(config)
	
}
