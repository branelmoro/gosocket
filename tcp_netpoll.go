package main

import (
    "fmt"
    "net"
    "os"
    "github.com/mailru/easygo/netpoll"
    "runtime"
    // "time"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
)

type abcd interface {
  bc()
}

type strct1 struct {
  net.Conn
  d int
}

// func (s *strct1) bc() {
//   fmt.Println("bc", s.bc)
// }

func main() {

    // var zz abcd

    // yz := strct1{}
    // yz.bc = func(){}

    // zz = yz

    // z1, ok := abc.(strct1)

    // fmt.Println("zz, z1, ok", zz, z1, ok)
    // return

    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

    conf := netpoll.Config{}
    conf.OnWaitError = func(err error) {
      fmt.Println("---OnWaitError err----", err)
    }
    poller, err := netpoll.New(&conf)


    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        abcdf := strct1{conn,2}
        fmt.Println("abcdf:", abcdf)

        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        // go handleRequest(conn)
        go handleConn(conn, poller)
    }
}


func handleConn(conn net.Conn, poller netpoll.Poller) {

  fmt.Println("goroutine count on connection - ", runtime.NumGoroutine())

  // Get netpoll descriptor with EventRead|EventEdgeTriggered.
  // desc := netpoll.Must(netpoll.HandleRead(conn))

  // evts := netpoll.EventOneShot | netpoll.EventPollerClosed | netpoll.EventErr | netpoll.EventWriteHup | netpoll.EventReadHup | netpoll.EventHup | netpoll.EventRead | netpoll.EventWrite | netpoll.EventEdgeTriggered

  // evts := netpoll.EventPollerClosed | netpoll.EventWriteHup | netpoll.EventReadHup | netpoll.EventHup | netpoll.EventRead | netpoll.EventWrite | netpoll.EventEdgeTriggered

  evts := netpoll.EventPollerClosed | netpoll.EventWriteHup | netpoll.EventReadHup | netpoll.EventHup | netpoll.EventRead | netpoll.EventEdgeTriggered

  desc := netpoll.Must(netpoll.Handle(conn, evts))

  // err := conn.SetReadDeadline(time.Now().Add(time.Second))
  // fmt.Println("---initial timeout err----", err)

  poller.Start(desc, func(ev netpoll.Event) {

    fmt.Println("Event ----------- ", ev)

    fmt.Println("goroutine count - ", runtime.NumGoroutine())

    switch {
      case ev&netpoll.EventRead != 0:
        handleRequest(conn)
      // case netpoll.EventWrite:
      // case netpoll.EventOneShot:
      // case netpoll.EventEdgeTriggered:
      // case netpoll.EventReadHup:
      // case netpoll.EventWriteHup:
      // case netpoll.EventHup:
      // case netpoll.EventErr:
      // case netpoll.EventPollerClosed:
      // default:
    }

  })

  conn.Write([]byte("Connection established."))

}


// Handles incoming requests.
func handleRequest(conn net.Conn) {

  buf := make([]byte, 3000)
  num_bytes, err := conn.Read(buf)
  fmt.Println("----after read---------------", string(buf[:num_bytes]), num_bytes, err)

  op := fmt.Sprintf("goroutine count - %d", runtime.NumGoroutine())

  conn.Write([]byte(op))
  // conn.Close()
}