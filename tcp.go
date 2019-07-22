package main

import (
    "fmt"
    "net"
    "os"
    "github.com/mailru/easygo/netpoll"
    "reflect"
    "runtime"
    "time"
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
        go handleRequest(conn)
        // go handleConn(conn)
    }
}



func handleConn(conn net.Conn) {


  fmt.Println("goroutine count on connection - ", runtime.NumGoroutine())



  // desc, err := netpoll.Handle(conn, netpoll.EventRead | netpoll.EventEdgeTriggered)
  // if err != nil {
  //   // handle error
  // }


  conf := netpoll.Config{}
  conf.OnWaitError = func(err error) {
    fmt.Println("---OnWaitError err----", err)
  }


  poller, err := netpoll.New(&conf)
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

  // k := 0
  err = conn.SetReadDeadline(time.Now().Add(time.Second))
  fmt.Println("---initial timeout err----", err)

  poller.Start(desc, func(ev netpoll.Event) {

    fmt.Println("Event ----------- ", ev)


    // k++

    // if k == 1 {
      handlePollEvent(conn, poller, desc)
    // }

    // poller.Stop(desc)
    // desc.Close()
    // conn.Close()

    fmt.Println("goroutine count - ", runtime.NumGoroutine())

    
  })


  conn.Write([]byte("Connection established."))

}


func handlePollEvent(conn net.Conn, poller netpoll.Poller, desc *netpoll.Desc) {
  // poller.Stop(desc)
  handleEvent(conn)
  // poller.Resume(desc)

}

func handleEvent(conn net.Conn) {

  var (
    err error
    num_bytes int
    buf []byte
  )


  fmt.Println("conn is ------------", reflect.TypeOf(conn))
  i :=0


  // err := conn.(*net.TCPConn).SetKeepAlive(false)
  // if err != nil {
  //   fmt.Println("----conn.SetKeepAlive failed--------")
  // }

  for {

    // conn.SetKeepAlivePeriod(time.Microsecond)


    // start message read

    buf = make([]byte, 3)
    fmt.Println("----before read----------i-----")
    num_bytes, err = conn.Read(buf)
    fmt.Println("----after read----------i-----")


    if err, ok := err.(net.Error); ok && err.Timeout() {
      fmt.Println("timeout occured================")
      // break
    }

    fmt.Println(err, buf, num_bytes)

    i += num_bytes


    fmt.Println(i%3, "----mod----------i-----", i)

    if i%3 == 0 {
      // message read done
      fmt.Println("mesage read done-----")
      err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
      fmt.Println("set SetReadDeadline limited-----err is------", err)
      time.Sleep(100 * time.Millisecond)
    } else {
      conn.SetReadDeadline( time.Time{})
      fmt.Println("set SetReadDeadline infinite----")
    }

    if i > 15 {
      fmt.Println("breaking after 15-----")
      break
    }
  }

}


// Handles incoming requests.
func handleRequest(conn net.Conn) {
  // Make a buffer to hold incoming data.
  // buf := make([]byte, 1024)
  // // Read the incoming connection into the buffer.
  // reqLen, err := conn.Read(buf)
  // fmt.Println(reqLen)
  // if err != nil {
  //   fmt.Println("Error reading:", err.Error())
  // }
  // // Send a response back to person contacting us.
  // conn.Write([]byte("Message received."))
  // // Close the connection when you're done with it.

  // handleEvent(conn)



  buf := make([]byte, 3000)
  num_bytes, err := conn.Read(buf)
  fmt.Println("----after read---------------", buf, num_bytes, err)

  fmt.Println(string(buf[:num_bytes]))



  conn.Close()
}