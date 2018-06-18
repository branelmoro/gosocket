package gosocket

import (
    // "fmt"
    // "net"
    // "os"
    // "github.com/mailru/easygo/netpoll"
    // "runtime"
)





type socketOpenCb func(*Conn, *[]byte) bool

type messageCb func(*Conn, *Message)

type errorCb func(*Conn)

var(
    OnWebsocketOpen socketOpenCb
    OnMessage messageCb
    OnError errorCb
    OnClose messageCb
)