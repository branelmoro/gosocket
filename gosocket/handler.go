package gosocket

// import (
//     "fmt"
//     "net"
//     "os"
//     "github.com/mailru/easygo/netpoll"
//     "runtime"
// )

type socketOpenCb func(Conn)

type messageCb func(Conn, []byte)

type errorCb func(Conn)

type closeCb func(Conn)


// var onWebsocketOpen socketOpenCb
// var onMessage messageCb
// var onError errorCb
// var onClose closeCb


// func OnWebsocketOpen(cb socketOpenCb) {
// 	onWebsocketOpen = cb
// }

// func OnMessage(cb messageCb) {
//     onMessage = cb
// }

// func OnError(cb errorCb) {
//     onError = cb
// }

// func OnClose(cb closeCb) {
//     onClose = cb
// }




var OnWebsocketOpen socketOpenCb
var OnMessage messageCb
var OnError errorCb
var OnClose closeCb