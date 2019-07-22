package gosocket

import (
    "fmt"
    "net"
    "os"
    // "time"
    "github.com/mailru/easygo/netpoll"
    // "runtime"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
)

func StartServer() {
    // fmt.Println(time.Second)
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
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        // go handleRequest(conn)
        // go handleConn(conn)

        poller, err1 := netpoll.New(nil)
        if err1 != nil {
            conn.Write([]byte("Unable to initialize netpoll... Closing Connection..."))
            conn.Close()
            return
        }

        // Get netpoll descriptor with EventRead|EventEdgeTriggered.
        desc := netpoll.Must(netpoll.Handle(conn, netpoll.EventRead | netpoll.EventEdgeTriggered))

        socketConn := Conn{conn: conn, desc: desc, poller: poller}
        
        go httpRequestHandler(&socketConn)
    }
}

