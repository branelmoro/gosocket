package main

import (
    "fmt"
    "net"
    // "os"
    // "github.com/mailru/easygo/netpoll"
    // "runtime"
    // "time"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = 3333
    CONN_TYPE = "tcp"
)

type wsInterface struct {
	OnWebsocketOpen func()
}


func main() {
	a := wsInterface{}

	fmt.Println(a)
	a.OnWebsocketOpen()

	conn, err := net.Dial(CONN_TYPE, fmt.Sprintf("%s:%d", CONN_HOST, CONN_PORT))
	if err != nil {
		// handle error
	}
	// fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	// status, err := bufio.NewReader(conn).ReadString('\n')

	fmt.Println(conn, err)

	data := []byte("hgsydhjfkdflma")

	conn.Write(data)

	handleRequest(conn)

	conn.Write(data)

	handleRequest(conn)

	conn.Write(data)

	handleRequest(conn)

 	conn.Close()
}


func handleRequest(conn net.Conn) {
	buf := make([]byte, 3000)
	num_bytes, err := conn.Read(buf)
	fmt.Println("----after read---------------", string(buf[:num_bytes]), num_bytes, err)
}