package gosocket

import (
	"fmt"
	"github.com/mailru/easygo/netpoll"
	"sync"
	"time"
)

const(
	// bit order
	// reserved, isConnClosed, pingReceived, pingSent, writeStatus, readStatus, closeSent, closeReceived

	// use logical OR `|` operation
	CLOSE_SENT          byte = 0x02 // 00000010
	CLOSE_RECEIVED      byte = 0x01 // 00000001

	PING_SENT           byte = 0x10 // 00010000
	PING_RECEIVED       byte = 0x20 // 00100000

	CONN_CLOSED         byte = 0x40 // 01000000

	WRITING_ON          byte = 0x08 // 00001000
	READING_ON          byte = 0x04 // 00000100

	// use logical AND `&` operation
	WRITING_OFF         byte = 0xf7 // 11110111
	READING_OFF         byte = 0xfb // 11111011
	PONG_SENT           byte = 0xdf // 11011111
	PONG_RECEIVED       byte = 0xef // 11101111
)


type wsConn struct {
	*Conn
	ConnData interface{}

	// local vars for conn
	_isClient bool
	connStatus byte
	pingSent []byte
	pingReceived []byte

	// locks for read and write
	_readLock sync.Mutex
	_writeLock sync.Mutex

	// lock for connection
	_connLock sync.Mutex

	flate *perMessageDeflate
}


func (ws *wsConn) closeConn() error {
	fmt.Println("calling closeConn------")
	ws._connLock.Lock()
    defer ws._connLock.Unlock()
    defer ws.server.delConn(ws)
    if ws.isConnClosed() {
    	return nil
    }
    err := ws.close()
    if err != nil {
    	ws.connStatus |= CONN_CLOSED
    }
    return err
}

func (ws *wsConn) setConnStatus(status byte) {
    ws._connLock.Lock()
    defer ws._connLock.Unlock()
    switch(status) {
    	case CLOSE_SENT, CLOSE_RECEIVED, WRITING_ON, READING_ON, PING_SENT, PING_RECEIVED:
    		ws.connStatus |= status
    		break
    	case WRITING_OFF, READING_OFF, PONG_SENT, PONG_RECEIVED:
    		ws.connStatus &= status
    		break
    }
}

func (ws *wsConn) handleTCPClose(err error) {
    ws._connLock.Lock()
    defer ws._connLock.Unlock()
    defer ws.server.delConn(ws)
    if ws.isConnClosed() {
    	return
    }
    ws.connStatus |= CONN_CLOSED
    writer := ws.writer()
    go OnError(writer, err)
	if !ws.isCloseReceived() && !ws.isCloseSent() {
		go OnClose(writer, NewCloseMsg(CC_ABNORMAL_CLOSE, ""))
	}
}

func (ws *wsConn) isCloseSent() bool {
	return ws.connStatus&CLOSE_SENT == CLOSE_SENT
}

func (ws *wsConn) isCloseReceived() bool {
	return ws.connStatus&CLOSE_RECEIVED == CLOSE_RECEIVED
}

func (ws *wsConn) isConnClosed() bool {
	return ws.connStatus&CONN_CLOSED == CONN_CLOSED
}

func (ws *wsConn) isPerMessageDeflateEnabled() bool {
	return ws.flate != nil
}

func (ws *wsConn) writer() *wsWriter {
	return &wsWriter{wsConn:ws}
}

func (ws *wsConn) reader() *wsReader {
	return &wsReader{
		wsConn: ws,
		isClosing: false,
	}
}

func (ws *wsConn) startReading() {
	ws.reader().start()
}

func openWebSocket(ws *wsConn) {

	// keep connection open forever
	ws.setReadTimeOut(time.Time{})

	writer := ws.writer()

	go OnWebsocketOpen(writer)

	// start netpoll
	// poller := *(ws.conn.poller)
	// desc := ws.conn.desc

	ws.fpoller().Start(ws.fdesc(), func(ev netpoll.Event) {

		// defer func() {
		// 	if r := recover(); r != nil {
		// 		fmt.Println("Recovered in f", r)
		// 		conn.Close()
		// 	}
		// }()

		fmt.Println("----------------------------------------------------------------",ev)

		ws.startReading()
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

	// add connection to server
	ws.server.addConn(ws)
}
