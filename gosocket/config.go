package gosocket

import "time"

// deflate settings
// {
//     "request_no_context_takeover": true,
//     "request_max_window_bits": 12,
//     "no_context_takeover": true,
//     "max_window_bits": 12,
//     "memory_level": 5
// }

type WsDeflateConf struct {
    ContextTakeover bool
    MaxWindowBits byte
    MemoryLevel byte
}

func NewWsDeflateConf() (*WsDeflateConf) {
    return &WsDeflateConf{
        ContextTakeover:    false,
        MaxWindowBits:      15,
        MemoryLevel:        8,
    }
}


type WsConf struct {

    WsDeflate *WsDeflateConf

    WsMaxFrameSize uint
    WsMaxMessageSize uint

    WsHeaderReadTimeout time.Duration
    WsMinByteRatePerSec uint
    WsCloseTimeout time.Duration

    OnWebsocketOpen func(WsWriter)
    OnMessage func(WsWriter, Message)
    OnText func(WsWriter, string)
    OnBinary func(WsWriter, []byte)
    OnError func(WsWriter, error)
    OnClose func(WsWriter, CloseMsg)
    OnPing func(WsWriter)
    OnPong func(WsWriter)
}

func NewWsConf() (*WsConf) {

    conf := &WsConf{
        WsMaxFrameSize:             65536,
        WsMaxMessageSize:           65536,

        WsHeaderReadTimeout:        1,
        WsMinByteRatePerSec:        100,
        WsCloseTimeout:             2,
    }

    conf.WsDeflate = NewWsDeflateConf()

    // on websocket connection open
    conf.OnWebsocketOpen = func(w WsWriter) {}

    // on message
    conf.OnMessage = func(w WsWriter, msg Message) {}

    // on message
    conf.OnText = func(w WsWriter, str string) {}

    // on message
    conf.OnBinary = func(w WsWriter, data []byte) {}

    // on error
    conf.OnError = func(w WsWriter, err error) {}

    // on connection close
    conf.OnClose = func(w WsWriter, msg CloseMsg) {}

    // on ping
    conf.OnPing = func(w WsWriter) {}

    // on pong
    conf.OnPong = func(w WsWriter) {}

    return conf

}




type ServerConf struct {
    Headers map[string]string

    Host string
    BindHosts []string
    Port uint16

    CertPrivate []byte
    CertPublic []byte    

    HttpHeaderTimeOut time.Duration
    HttpMaxHeaderSize uint

    NetworkBandWidth uint
    MaxWsConnections uint

    OnMalformedRequest func(HttpRequest)
    OnHttpRequest func(HttpWriter, HttpRequest)
    OnWebsocketRequest func(HttpWriter, HttpRequest)

    WsConf *WsConf
}

func NewServerConf() (*ServerConf) {
    conf := &ServerConf{
        Host:                       "localhost",
        BindHosts:                  []string{"127.0.0.1"},
        Port:                       3333,

        HttpHeaderTimeOut:          20,
        HttpMaxHeaderSize:          8192,

        NetworkBandWidth:           0,
        MaxWsConnections:            0,

        OnMalformedRequest: func(r HttpRequest) {},
        OnHttpRequest: func(w HttpWriter, r HttpRequest) {},
        OnWebsocketRequest: func(w HttpWriter, r HttpRequest) {},
    }

    conf.WsConf = NewWsConf()
    return conf
}




type ClientConf struct {
    Headers map[string]string
    WsConf *WsConf
}

func NewClientConf() (*ClientConf) {
    conf := &ClientConf{}

    conf.WsConf = NewWsConf()
    return conf
}
