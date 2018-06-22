package gosocket

import (
	"fmt"
	"net"
	"io"
	"time"
	// "bufio"
	"github.com/mailru/easygo/netpoll"
)

// Conn represents single connection instance.
type Conn struct {
	conn	 net.Conn
	desc 	 *netpoll.Desc
	poller   *netpoll.Poller
	message  []byte
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

func readRequestTillBlankNewline(c *net.Conn, nl_count int, byte_size int) ([]byte, int){

	conn := *c

	var (
		read_bytes []byte
		prev_byte byte
		// err string
	)

	buff := make([]byte, 1)
	// Read the incoming connection into the buffer.

	prev_byte = 0

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

	return read_bytes, byte_size

}


// Handles incoming requests.
func handleConnection(conn net.Conn) {

	var (
		read_bytes []byte
		headers map[string]string
		sec_web_accept string
	)

	request_type := "tcp"

	timeoutDuration := 10 * time.Second
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	req_bytes, byte_size := readRequestTillBlankNewline(&conn, 2, 0)

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

	poller, err1 := netpoll.New(nil)
	if err1 != nil {
		conn.Write([]byte("Unable to initialize netpoll... Closing Connection..."))
		conn.Close()
		return
	}

	// Get netpoll descriptor with EventRead|EventEdgeTriggered.
	desc := netpoll.Must(netpoll.Handle(conn, netpoll.EventRead | netpoll.EventEdgeTriggered))

	connection := Conn{conn: conn, desc: desc, poller: &poller}

	
	// validate content in read_bytes
    req, req_len, err2 := getRequest(&read_bytes)

    fmt.Println(req, req_len, err2)

	if err2 == "" {
		headers, err3 := getHeaders(&read_bytes,(req_len+1))
		fmt.Println(headers)
		if err3 == "" {
			request_type = "http"
			if v1, v2 := headers["Upgrade"], headers["Sec-WebSocket-Key"]; v1 != "" && v1 == "websocket" && v2 != "" && req[0] == "GET" {
				request_type = "websocket"
				sec_web_accept = getSecWebSocketAccept(headers["Sec-WebSocket-Key"])
			}
		}
	}

	fmt.Println(headers)

	switch(request_type) {
		case "websocket":
			is_valid := OnWebsocketOpen(&connection, &read_bytes)

			// handle error
			// OnError(connection)

			// upgrade to websocket connection
			if is_valid {
				fmt.Println(headers["Sec-WebSocket-Key"], sec_web_accept)
				resp := append([]byte("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "), []byte(sec_web_accept)...)
				// resp := append([]byte("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nContent-Encoding: identity\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "), []byte(sec_web_accept)...)

				resp = append(resp, []byte("\r\n\r\n")...)
				fmt.Println(string(resp), resp)
				connection.conn.Write(resp)
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


		// defer func() {
		// 	if r := recover(); r != nil {
		// 		fmt.Println("Recovered in f", r)
		// 		conn.Close()
		// 	}
		// }()


		fmt.Println(ev)

		// OnMessage(c, conn.Read())
		// fmt.Println(conn.Read())

		conn.readMessages()
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


func closeWebsocketPanic(c *Conn) {
	conn := *c
	if r := recover(); r != nil {
		fmt.Println("Recovered in f", r)
		conn.Close()
	}
}