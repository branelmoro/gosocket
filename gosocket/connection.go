package gosocket

import (
	"fmt"
    "net"
    "github.com/mailru/easygo/netpoll"
)

// Conn represents single connection instance.
type Conn struct {
	conn     net.Conn
	desc 	 *netpoll.Desc
	poller   *netpoll.Poller
}

func (c *Conn) Read() []byte {
	b := make([]byte, 1000)
	x, err := c.conn.Read(b)
	fmt.Println("message reading:", x, err, string(b),)
	return b
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


// Handles incoming requests.
func handleConnection(conn net.Conn) {
    
	poller, err := netpoll.New(nil)
	if err != nil {
		// handle error
	} else {

		// Get netpoll descriptor with EventRead|EventEdgeTriggered.
		desc := netpoll.Must(netpoll.Handle(conn, netpoll.EventRead | netpoll.EventEdgeTriggered))

		poller.Start(desc, func(ev netpoll.Event) {
			fmt.Println("start-------------------")
			fmt.Println(ev)
			b := make([]byte, 1000)
			x, err := conn.Read(b)
			fmt.Println("fisrt start reading:", x, err, string(b))
			poller.Stop(desc)

			OnWebsocketOpen(conn)

			go upgrateToWebSocket(conn)
		})
	}
}

func upgrateToWebSocket(conn net.Conn) {

	poller, err := netpoll.New(nil)
	if err != nil {
		// handle error

		OnError(conn)
	} else {

		// Get netpoll descriptor with EventRead|EventEdgeTriggered.
		desc := netpoll.Must(netpoll.Handle(conn, netpoll.EventRead | netpoll.EventEdgeTriggered))

		connection := Conn{conn: conn, desc: desc, poller: &poller}

		poller.Start(desc, func(ev netpoll.Event) {

			fmt.Println("new OnMessage call-------------------")
			fmt.Println(ev)

			OnMessage(connection, connection.Read())
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
}
