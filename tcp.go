package main

import (
    "fmt"
    "net"
    "os"
    "github.com/mailru/easygo/netpoll"
    "runtime"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
)

func main() {
    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        // go handleRequest(conn)
        go handleConn(conn)
    }
}



func handleConn(conn net.Conn) {


  fmt.Println("goroutine count on connection - ", runtime.NumGoroutine())



  // desc, err := netpoll.Handle(conn, netpoll.EventRead | netpoll.EventEdgeTriggered)
  // if err != nil {
  //   // handle error
  // }


  poller, err := netpoll.New(nil)
  if err != nil {
    // handle error
  }

  // Get netpoll descriptor with EventRead|EventEdgeTriggered.
  // desc := netpoll.Must(netpoll.HandleRead(conn))
  desc := netpoll.Must(netpoll.Handle(conn, netpoll.EventRead | netpoll.EventWrite | netpoll.EventEdgeTriggered))


  // netpoll.CallbackFn = func (ev netpoll.Event) {
  //   fmt.Println("----------")
  //   fmt.Println(ev)
  // }

  poller.Start(desc, func(ev netpoll.Event) {

    fmt.Println("Event ----------- ", ev)

    buf := make([]byte, 1024)
    data, err := conn.Read(buf)
    fmt.Println(ev, err, string(buf), data)

    // poller.Stop(desc)
    // desc.Close()
    // conn.Close()

    fmt.Println("goroutine count - ", runtime.NumGoroutine())

    
  })


  conn.Write([]byte("Connection established."))

}




// Handles incoming requests.
func handleRequest(conn net.Conn) {
  // Make a buffer to hold incoming data.
  buf := make([]byte, 1024)
  // Read the incoming connection into the buffer.
  reqLen, err := conn.Read(buf)
  fmt.Println(reqLen)
  if err != nil {
    fmt.Println("Error reading:", err.Error())
  }
  // Send a response back to person contacting us.
  conn.Write([]byte("Message received."))
  // Close the connection when you're done with it.
  conn.Close()
}