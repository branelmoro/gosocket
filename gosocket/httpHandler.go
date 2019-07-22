package gosocket

import (
	// "fmt"
	// "net"
	// "io"
	// "time"
)

func httpRequestHandler(conn *Conn) {
	reader := &httpReader{
		Conn: conn,
	}
	req, err := reader.readRequest()
	if err == nil {
		if req.isWebSocketRequest() {
			OnWebsocketRequest(&httpWriter{Conn:conn, req: req}, req)
		} else {
			OnHttpRequest(&httpWriter{Conn:conn, req: req}, req)
			conn.close()
		}
	} else {
		validAdminRequest := false
		if validAdminRequest {
			// process
		} else {
			conn.close()
			OnMalformedRequest(req)
		}
	}
}
