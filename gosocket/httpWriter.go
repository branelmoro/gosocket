package gosocket

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net"
	"strings"
	"time"
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

type httpWriter struct {
	*Conn
	req *httpRequest
}

func (w *httpWriter) Write(data []byte) error {
	var(
		numBytes int
		cntBytes int
		err error
	)
	sec := 10
	minBytes := sec * w.server.wsMinByteRatePerSec
	err = w.setWriteTimeOut(time.Now().Add(time.Duration(sec) * time.Second))
	if err != nil {
		return newSetWriteTimeoutError(err)
	}
	startIndex := 0
	cntBytes = 0
	for startIndex != len(data) {
		numBytes, err = w.write(data[startIndex:])
		cntBytes += numBytes
		if e, ok := err.(net.Error); ok && e.Timeout() {
			// timeout occured
			if cntBytes < minBytes {
				// return error, connection accept less data (numbytes bytes data) for 10 second
				// expecting minBytes (w.server.wsMinByteRatePerSec per second)
				return newSlowDataWriteError(cntBytes, sec)
			}
			err = w.setWriteTimeOut(time.Now().Add(time.Duration(sec) * time.Second))
			if err != nil {
				return newSetWriteTimeoutError(err)
			}
			err = nil
			cntBytes = 0
		}
		if err != nil{
			break
		}
		startIndex += numBytes
	}
	return err
}

func (w *httpWriter) getSecWebSocketAccept() string {
	str := append([]byte(w.req.header["sec-websocket-key"]), []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")...)
    // s := "dGhlIHNhbXBsZSBub25jZQ==258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
    h := sha1.New()
    h.Write(str)
    bs := h.Sum(nil)
    return base64.StdEncoding.EncodeToString(bs)
}

func (w *httpWriter) UpgradeToWebsocket(options *WsOptions) error {
	var(
		flate *perMessageDeflate
		err error
	)
	if !w.req.isWebSocketRequest() {
		return err
	}

	header := make(map[string]string)

	if options != nil && options.Headers != nil {
		// add extra headers

	}

	header["upgrade"] = "websocket"
	header["connection"] = "Upgrade"
	header["Sec-WebSocket-Accept"] = w.getSecWebSocketAccept()
	header["Sec-WebSocket-Version"] = "13"

	if true {
		if ext, ok := w.req.header["sec-websocket-extensions"]; ok && ext == "permessage-deflate" {
			header["sec-websocket-extensions"] = "permessage-deflate"
			flate, err = newPerMessageDeflate(9)
		}
	}

	bytes := []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", HttpSwitchingProtocols, httpStatusText[HttpSwitchingProtocols]))
	for key, val := range header {
		bytes = append(bytes, []byte(strings.Title(key) + ": " + val + "\r\n")...)
	}
	// add \r\n. end of header and http protocol
	bytes = append(bytes, 0xd, 0xa)
	err = w.Write(bytes)
	fmt.Println(err, string(bytes))
	if err == nil {
		conn := wsConn{
			Conn: w.Conn,
			// ConnData: options.WsData,
			_isClient: true,
			flate: flate,
		}
		openWebSocket(&conn)
	} else {
		// close connection
		w.close()
	}
	return err
}

func (w *httpWriter) Close() error {
	return w.close()
}
