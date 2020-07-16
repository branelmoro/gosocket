package gosocket

import (
	"fmt"
	"net"
	"time"
    "github.com/mailru/easygo/netpoll"
)

type WsOptions struct {
	Headers map[string]string
	WsData interface{}
}

type HttpWriter interface {
	Write([]byte) error
	UpgradeToWebsocket(*WsOptions) error
	Close() error
}

type httpWriterRequest interface {
	Header() map[string]string
	isWebSocketRequest() bool
}

type httpWriter struct {
	conn httpWriterConn
	req httpWriterRequest
	wsH *wsHandler
}

func (w *httpWriter) Write(data []byte) error {
	var(
		numBytes int
		cntBytes int
		err error
	)
	sec := 10
	minBytes := sec * w.wsH.minByteRate
	err = w.conn.setWriteTimeOut(time.Now().Add(time.Duration(sec) * time.Second))
	if err != nil {
		return newSetWriteTimeoutError(err)
	}
	startIndex := 0
	cntBytes = 0
	for startIndex != len(data) {
		// fmt.Println(startIndex)
		// fmt.Println(string(data[startIndex:]), startIndex)
		numBytes, err = w.conn.write(data[startIndex:])
		cntBytes += numBytes
		if e, ok := err.(net.Error); ok && e.Timeout() {
			// timeout occured
			if cntBytes < minBytes {
				// return error, connection accept less data (numbytes bytes data) for 10 second
				// expecting minBytes (w.conn.minSpeed() per second)
				return newSlowDataWriteError(cntBytes, sec)
			}
			err = w.conn.setWriteTimeOut(time.Now().Add(time.Duration(sec) * time.Second))
			if err != nil {
				return newSetWriteTimeoutError(err)
			}
			err = nil
			cntBytes = 0
		}
		if err != nil {
			return newWriteError(err)
		}
		startIndex += numBytes
	}
	return err
}

func (w *httpWriter) UpgradeToWebsocket(options *WsOptions) error {
	var(
		flate *perMessageDeflate
		err error
	)
	if !w.req.isWebSocketRequest() {
		return fmt.Errorf("Request is not websocket request... Can't upgrate to websocket...")
	}

	requestHeader := w.req.Header()

	res := httpResponse {
		protocol: "HTTP/1.1",
		code: HttpSwitchingProtocols,
		headers: generateWsUpgradeHeader(requestHeader, options, w.wsH.deflateConf),
	}

	bytes := res.toBytes()

	err = w.Write(bytes)
	fmt.Println(err, string(bytes))
	if err == nil {
		if _, ok := res.headers["sec-websocket-extensions"]; ok {
			flate, err = newPerMessageDeflate(9)
		}

		if c, ok := w.conn.(net.Conn); ok {

	        // evts := netpoll.EventOneShot | netpoll.EventPollerClosed | netpoll.EventErr | netpoll.EventWriteHup | netpoll.EventReadHup | netpoll.EventHup | netpoll.EventRead | netpoll.EventWrite | netpoll.EventEdgeTriggered
	        evts := netpoll.EventPollerClosed | netpoll.EventWriteHup | netpoll.EventReadHup | netpoll.EventHup | netpoll.EventRead | netpoll.EventEdgeTriggered

	        // Get netpoll descriptor with EventRead|EventEdgeTriggered.
	        desc := netpoll.Must(netpoll.Handle(c, evts))

	        poller := w.wsH.wsSH.poller

	        sConn := &serverConn{ &conn{ c, desc, poller, w.wsH } }

			ws := &wsConn{
				netpollConn: sConn,
				// ConnData: options.WsData,
				flate: flate,
			}

			go openWebSocket(ws, poller, desc)



			// if err == nil {
			// 	// add connection to wsHandler
			// 	sConn.handler.wsSH.addConn(ws)
			// } else {
			// 	go writer.wsH().onError(writer, err)

			// 	// close connection
			// 	msg := NewCloseMsg(CC_UNEXPECTED_ERROR, "poller start failed")
			// 	writer.Close(msg)
			// }
		} else {

		}
	} else {
		// close connection
		w.conn.close()
	}
	return err
}

func (w *httpWriter) Close() error {
	return w.conn.close()
}
