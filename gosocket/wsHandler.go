package gosocket

import(
	"time"
)

type wsHandler struct {
    isMaskRequired bool

    maxFrameSize int
    maxMessageSize int

    headerReadTimeout time.Duration
    closeTimeout time.Duration

    maxByteRate int
    minByteRate int

	onWebsocketOpen func(WsWriter)
	onMessage func(WsWriter, Message)
	onText func(WsWriter, string)
	onBinary func(WsWriter, []byte)
	onError func(WsWriter, error)
	onClose func(WsWriter, CloseMsg)
	onPing func(WsWriter)
	onPong func(WsWriter)

	deflateConf *WsDeflateConf

	wsSH *wsServerHandler
}

func (wc *wsHandler) setDefault(conf *WsConf) {

	if conf != nil && conf.WsMaxFrameSize > 0 {
		wc.maxFrameSize = int(conf.WsMaxFrameSize)
	} else {
		wc.maxFrameSize = 65536
	}

	if conf != nil && conf.WsMaxMessageSize > 0 {
		wc.maxMessageSize = int(conf.WsMaxMessageSize)
	} else {
		wc.maxMessageSize = 65536
	}

	if conf != nil && conf.WsHeaderReadTimeout > 0 {
		wc.headerReadTimeout = conf.WsHeaderReadTimeout
	} else {
		wc.headerReadTimeout = 1
	}

	if conf != nil && conf.WsMinByteRatePerSec > 0 {
		wc.minByteRate = int(conf.WsMinByteRatePerSec)
	} else {
		wc.minByteRate = 100
	}
	wc.maxByteRate = 0

	if conf != nil && conf.WsCloseTimeout > 0 {
		wc.closeTimeout = conf.WsCloseTimeout
	} else {
		wc.closeTimeout = 2
	}

	// connection handlers
	if conf != nil && conf.OnWebsocketOpen != nil {
		wc.onWebsocketOpen = conf.OnWebsocketOpen
	} else {
		wc.onWebsocketOpen = func(w WsWriter) {}
	}

	if conf != nil && conf.OnMessage != nil {
		wc.onMessage = conf.OnMessage
	} else {
		wc.onMessage = func(w WsWriter, m Message) {}
	}

	if conf != nil && conf.OnText != nil {
		wc.onText = conf.OnText
	} else {
		wc.onText = func(w WsWriter, s string) {}
	}

	if conf != nil && conf.OnBinary != nil {
		wc.onBinary = conf.OnBinary
	} else {
		wc.onBinary = func(w WsWriter, b []byte) {}
	}

	if conf != nil && conf.OnError != nil {
		wc.onError = conf.OnError
	} else {
		wc.onError = func(w WsWriter, e error) {}
	}

	if conf != nil && conf.OnClose != nil {
		wc.onClose = conf.OnClose
	} else {
		wc.onClose = func(w WsWriter, c CloseMsg) {}
	}

	if conf != nil && conf.OnPing != nil {
		wc.onPing = conf.OnPing
	} else {
		wc.onPing = func(w WsWriter) {}
	}

	if conf != nil && conf.OnPong != nil {
		wc.onPong = conf.OnPong
	} else {
		wc.onPong = func(w WsWriter) {}
	}
}
