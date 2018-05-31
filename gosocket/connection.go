package gosocket

import (
    "net"
    "github.com/mailru/easygo/netpoll"
)

// Conn represents single connection instance.
type Conn struct {
	conn     net.Conn
	desc     netpoll.Desc
	poller   netpoll.Poller
}

func (c *Conn) Read() string {
	b := make([]byte, 1000)
	for {
		x, err := c.Conn.Read(b)
		fmt.Println("start reading:", x, err, string(b), b)
		if err != nil {
		  break
		}
		if string(b) == "" {
		  break
		}
	}
}

func (c *Conn) Write() {

}

func (c *Conn) Close() {
	conn.poller.close()
	conn.desc.close()
	conn.Conn.close()
}


func upgrateToWebSocket(conn *net.Conn) {



}


// Handles incoming requests.
func handleConnection(conn *net.Conn) {
    
	poller, err := netpoll.New(nil)
	if err != nil {
		// handle error
	} else {

		// Get netpoll descriptor with EventRead|EventEdgeTriggered.
		desc := netpoll.Must(netpoll.Handle(conn, netpoll.EventRead | netpoll.EventEdgeTriggered))

		connection := Conn{conn: conn, desc: desc, poller: poller}

		poller.Start(desc, func(ev netpoll.Event) {



			fmt.Println("start-------------------")
			fmt.Println(ev)


			connection.Read()
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
		}
	}
}