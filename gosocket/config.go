package gosocket

import "time"

type ServerConf struct {
    Host string
    BindHosts []string
    Port uint16

    CertPrivate []byte
    CertPublic []byte    

    HttpRquestTimeOut time.Duration
    HttpMaxRequestLineSize uint
    HttpMaxHeaderSize uint

    WsMaxFrameSize uint
    WsMaxMessageSize uint

    WsHeaderReadTimeout time.Duration
    WsMinByteRatePerSec uint
    WsCloseTimeout time.Duration

    NetworkBandWidth uint
    MaxWsConnection uint
}


func NewConf() (*ServerConf) {
    return &ServerConf{
        Host:                       "localhost",
        BindHosts:                  []string{"127.0.0.1"},
        Port:                       3333,

        HttpRquestTimeOut:          20,
        HttpMaxRequestLineSize:     1024,
        HttpMaxHeaderSize:          8192,

        WsMaxFrameSize:             65536,
        WsMaxMessageSize:           65536,

        WsHeaderReadTimeout:        1,
        WsMinByteRatePerSec:        100,
        WsCloseTimeout:             2,

        NetworkBandWidth:           0,
        MaxWsConnection:            0,
    }
}