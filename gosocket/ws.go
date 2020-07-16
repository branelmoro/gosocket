package gosocket

import (
	"fmt"
	"github.com/mailru/easygo/netpoll"
	"io"
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
	netpollConn
	ConnData interface{}

	// local vars for conn
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
    if ws.wsH().wsSH != nil {
    	defer ws.wsH().wsSH.delConn(ws)
    }
    if ws.isConnClosed() {
    	return nil
    }
    err := ws.close()
    // if err != nil {
    	ws.connStatus |= CONN_CLOSED
    // }
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
    if ws.wsH().wsSH != nil {
    	defer ws.wsH().wsSH.delConn(ws)
    }
    if ws.isConnClosed() {
    	return
    }
    ws.stopPoller()
    ws.connStatus |= CONN_CLOSED
    writer := ws.writer()
    if !ws.isCloseReceived() || !ws.isCloseSent() {
    	go ws.wsH().onClose(writer, NewCloseMsg(CC_ABNORMAL_CLOSE, ""))
	    go ws.wsH().onError(writer, err)
	}
	// if !ws.isCloseReceived() && !ws.isCloseSent() {
	// 	go ws.wsH().onClose(writer, NewCloseMsg(CC_ABNORMAL_CLOSE, ""))
	// }
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

func openWebSocket(ws *wsConn, poller netpoll.Poller, desc *netpoll.Desc) {

	// keep connection open forever
	ws.setReadTimeOut(time.Time{})

	writer := ws.writer()

	go OnWebsocketOpen(writer)

	err := poller.Start(desc, func(ev netpoll.Event) {

		fmt.Println("----------------------------------------------------------------",ev)

		switch {
			case ev&netpoll.EventPollerClosed != 0:
				// poller is closed assign connection to another poller
				break
			case ev&netpoll.EventHup != 0:
				// connection is closed by client
				ws.handleTCPClose(newTCPError(ERR_TCP_CLOSE, io.EOF))
				break
			case ev&netpoll.EventReadHup != 0:
				// connection stopped reading and closed by client
				ws.handleTCPClose(newTCPError(ERR_TCP_CLOSE, io.EOF))
				break
			case ev&netpoll.EventWriteHup != 0:
				// connection stopped writing and closed by client
				ws.handleTCPClose(newTCPError(ERR_TCP_CLOSE, io.EOF))
				break
			case ev&netpoll.EventRead != 0:
				ws.startReading()
				break
			case ev&netpoll.EventWrite != 0:
				break
			case ev&netpoll.EventOneShot != 0:
				break
			case ev&netpoll.EventEdgeTriggered != 0:
				break
			case ev&netpoll.EventErr != 0:
				break
			default:
				break
		}
	});

	if err == nil {
		// add connection to wsHandler
		if ws.wsH().wsSH != nil {
			ws.wsH().wsSH.addConn(ws)
		}
	} else {
		// close connection
		msg := NewCloseMsg(CC_UNEXPECTED_ERROR, "poller start failed")
		go writer.Close(msg)
		go ws.wsH().onError(writer, err)
	}
}
