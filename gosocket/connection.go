package gosocket

import (
	"fmt"
	"net"
	"io"
	"time"
	// "bufio"
	"github.com/mailru/easygo/netpoll"
	"strings"
)

// Conn represents single connection instance.
type Conn struct {
	conn	 net.Conn
	desc 	 *netpoll.Desc
	poller   *netpoll.Poller
}

func (c *Conn) Read() *[]byte {

	var read_bytes []byte

	buff_size := 1

	timeoutDuration := 1 * time.Millisecond
	fmt.Println("Time-----", timeoutDuration)
	c.conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	buff := make([]byte, buff_size)
	for {
		num_bytes, err := c.conn.Read(buff)
		// fmt.Println("Bytes received:", num_bytes, err, string(buff), time.Now())

		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}

		read_bytes = append(read_bytes, buff...)

		if num_bytes < buff_size {
			break
		} else {
			c.conn.SetReadDeadline(time.Now().Add(timeoutDuration))
		}
	}
	return &read_bytes
}

func (c *Conn) Write(data []byte) {
	c.conn.Write(data)
}

func (c *Conn) Close() {
	p := *c.poller
	p.Stop(c.desc)
	c.desc.Close()
	c.conn.Close()
}

func readRequestTillBlankNewline(c *net.Conn, nl_count int, byte_size int) ([]byte, int, string){

	conn := *c

	var read_bytes []byte

	buff := make([]byte, 1)
	// Read the incoming connection into the buffer.

	prev_byte := 0

	cnt_nl := 0

	for {
		num_bytes, err := conn.Read(buff)

		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}

		byte_size += num_bytes
		read_bytes = append(read_bytes, buff...)

		// logic to detect carriage return - 13 and newline - 10 characters
		if prev_byte == 13 && buff[0] == 10 {
			cnt_nl += 1
			if cnt_nl == nl_count {
				break
			}
		} else if prev_byte != 10 || buff[0] != 13 {
			cnt_nl = 0
		}
		prev_byte = buff[0]

		if num_bytes == 0 {
			break
		}
		conn.SetReadDeadline(time.Now().Add(time.Millisecond))

		if byte_size > 8192 {
			break
		}
	}

	return read_bytes, byte_size, err

}


// Handles incoming requests.
func handleConnection(conn net.Conn) {

	var read_bytes []byte

	request_type := "tcp"

	timeoutDuration := 10 * time.Second
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	req_bytes, byte_size, err := readRequestTillNewlines(&conn, 2, 0)

	read_bytes = append(read_bytes, req_bytes...)

	if byte_size == 0 {
		conn.Write([]byte("No data received till 10 seconds...Closing Connection..."))
		conn.Close()
		return
	}

	if byte_size > 8192 {
		conn.Write([]byte("Request/Header size more than available buffer - 8192 bytes...Closing Connection..."))
		conn.Close()
		return
	}


// 	if len(read_bytes) == 0 {

// 		// Send a response back to person contacting us.
// 		conn.Write([]byte(`HTTP/1.1 200 OK
// Server: Apache/2.2.14 (Win32)
// ETag: "10000000565a5-2c-3e94b66c2e680"
// Accept-Ranges: bytes
// Connection: close
// Content-Type: text/html
// X-Pad: avoid browser bug

// <html><body><h1>No Request data received!</h1></body></html>
// 		`))
// 		// Close the connection when you're done with it.
// 		conn.Close()
// 		return
// 	}
	
	// validate content in read_bytes

	poller, err := netpoll.New(nil)
	if err != nil {
		conn.Write([]byte("Unable to initialize netpoll... Closing Connection..."))
		conn.Close()
		return
	}

	// Get netpoll descriptor with EventRead|EventEdgeTriggered.
	desc := netpoll.Must(netpoll.Handle(conn, netpoll.EventRead | netpoll.EventEdgeTriggered))

	connection := Conn{conn: conn, desc: desc, poller: &poller}


    req, req_len, err := getRequest(&read_bytes)

	if err == "" {
		headers, err :=getHeaders(&a,(req_len+1))
		if err == "" {
			request_type = "http"
			if v1, v2 := headers["Upgrade"], headers["Sec-WebSocket-Key"]; v1 != "" && v1 == "websocket" && v2 != "" && req[0] == "GET" {
				request_type = "websocket"
			}
		}
	}

	switch(request_type) {
		case "websocket":
			is_valid := OnWebsocketOpen(&connection, &read_bytes)

			// handle error
			// OnError(connection)

			// upgrade to websocket connection
			if is_valid {
				upgrateToWebSocket(&connection)
			}
			break
		case "http":
		case "tcp":
			connection.Write([]byte("http or tcp not allowed... closing connection.."))
			connection.Close()
			break
	}

}


func upgrateToWebSocket(c *Conn) {

	conn := *c

	poller := *conn.poller
	desc := conn.desc

	poller.Start(desc, func(ev netpoll.Event) {

		fmt.Println(ev)

		OnMessage(c, conn.Read())
		// if ev&netpoll.EventReadHup != 0 {
		//   // poller.Stop(desc)
		//   conn.Close()
		//   return
		// }

		// hr, err := ioutil.ReadAll(conn)
		// fmt.Println(hr)
		// if err != nil {
		//   // handle error
		//
	})

}
