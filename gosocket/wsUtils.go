package gosocket

import (
	"crypto/sha1"
	"encoding/base64"
	// "fmt"
	// "net"
	// "strings"
	// "time"
)

func generateSecWebSocketAccept(secWebsocketKey string) string {
    str := append([]byte(secWebsocketKey), []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")...)
    h := sha1.New()
    h.Write(str)
    bs := h.Sum(nil)
    return base64.StdEncoding.EncodeToString(bs)
}

func generateWsUpgradeHeader(requestHeader map[string]string, options *WsOptions, deflateConf *WsDeflateConf) map[string]string {

	headers := make(map[string]string)

	if options != nil && options.Headers != nil {
		// add extra headers
		headers = options.Headers
	}

	headers["upgrade"] = "websocket"
	headers["connection"] = "Upgrade"
	headers["Sec-WebSocket-Accept"] = generateSecWebSocketAccept(requestHeader["sec-websocket-key"])
	headers["Sec-WebSocket-Version"] = "13"

	if deflateConf != nil {
		if ext, ok := requestHeader["sec-websocket-extensions"]; ok && ext == "permessage-deflate" {
			headers["sec-websocket-extensions"] = "permessage-deflate"
		}
	}

	return headers
}