package gosocket

import (
    "crypto/rand"
    "fmt"
    "crypto/tls"
    "net"
    // "sync"
    "time"
    "github.com/mailru/easygo/netpoll"
)

type server struct {
    listener net.Listener

    host string
    bindHosts []string
    port int

    certPrivate []byte
    certPublic []byte

    isListenerOn bool

    httpHeaderTimeOut time.Duration
    httpMaxHeaderSize int

    maxConnections int
    networkBandWidth int

    minByteRatePerSec int

    wsH *wsHandler

    wsSH *wsServerHandler

    onMalformedRequest func(HttpRequest)
    onHttpRequest func(HttpWriter, HttpRequest)
    onWebsocketRequest func(HttpWriter, HttpRequest)

}


func (s *server) startRateLimiter() {
    bw := s.networkBandWidth * 1024 * 1024
    for {
        iOps := s.wsSH.cntReadOps + s.wsSH.cntWriteOps
        if bw == 0 {
            s.wsH.maxByteRate = bw
        } else {
            s.wsH.maxByteRate = bw/int(iOps)
        }
        if s.minByteRatePerSec < s.wsH.maxByteRate {
            s.wsH.minByteRate = s.minByteRatePerSec
        } else {
            s.wsH.minByteRate = s.wsH.maxByteRate
        }
        time.Sleep(time.Second)
    }
}

func (s *server) startListener() error {
    var(
        listener net.Listener
        err error
    )

    network := "tcp"
    address := fmt.Sprintf("%s:%d", s.host, s.port)

    if len(s.certPublic) > 0 || len(s.certPrivate) > 0 {
        cert, err := tls.X509KeyPair(s.certPublic, s.certPrivate)
        if err != nil {
            return err
        }
        config := tls.Config{Certificates: []tls.Certificate{cert}}
        config.Rand = rand.Reader
        listener, err = tls.Listen(network, address, &config)
    } else {
        listener, err = net.Listen(network, address)
    }

    if err != nil {
        return err
    }
    s.listener = listener
    fmt.Println(fmt.Sprintf("Listening on %s network address %s", s.listener.Addr().Network(), s.listener.Addr().String()))
    return nil
}

func (s *server) logError(err error) error {
    fmt.Println(err.Error())
    return err
}

func (s *server) Run() {
    var(
        err error
        conn net.Conn
    )
    err = s.startListener()
    if err != nil {
        s.logError(err)
        fmt.Println("Error in starting listener: ", err)
        return
    }

    s.wsSH.poller, err = netpoll.New(nil)
    if err != nil {
        s.logError(err)
        fmt.Println("Error in creating connection poller: ", err)
        return
    }

    if s.networkBandWidth > 0 {
        go s.startRateLimiter()
    }

    // mark server listening
    s.isListenerOn = true

    for {
        // wait if connection limit is reached
        for s.isListenerOn && s.maxConnections > 0 && s.wsSH.wsCount() >= s.maxConnections {
            fmt.Println("Max connection limit reached, waiting for one second--- limit", s.maxConnections)
            time.Sleep(time.Second)
        }

        // Listen for an incoming connection.
        conn, err = s.listener.Accept()
        if err != nil {
            if s.isListenerOn {
                s.logError(err)
            }
            break
        }

        go s.handleConnection(conn)
    }

    // close all opened connections
    s.wsSH.shutDown()

}

func (s *server) handleConnection(conn net.Conn) {
    hConn := &httpconn{conn}
    reader := &httpRequestReader{&httpReader{ conn: hConn, maxHeaderSize: s.httpMaxHeaderSize }, s.host}
    req, err := reader.readRequest()
    if err == nil {
        w := &httpWriter{conn:hConn, req: req, wsH: s.wsH}
        if req.isWebSocketRequest() {
            s.onWebsocketRequest(w, req)
        } else {
            s.onHttpRequest(w, req)
            conn.Close()
        }
    } else {
        conn.Close()
        s.onMalformedRequest(req)
    }
}

func (s *server) Shutdown() {
    err := s.listener.Close()
    if err != nil {
        s.logError(err)
        fmt.Println("Error Server Shutdown: ", err)
    }

    s.isListenerOn = false
}

func (s *server) Restart() {
    s.Shutdown()
    s.Run()
}

func (s *server) Status() {


}

func NewServer(conf *ServerConf) Server {

    if conf == nil {
        conf = NewServerConf()
    }

    s := server{}

    s.host = conf.Host
    s.bindHosts = conf.BindHosts
    s.port = int(conf.Port)

    s.certPrivate = conf.CertPrivate
    s.certPublic = conf.CertPublic

    s.httpHeaderTimeOut = conf.HttpHeaderTimeOut
    s.httpMaxHeaderSize = int(conf.HttpMaxHeaderSize)

    s.networkBandWidth = int(conf.NetworkBandWidth)
    s.maxConnections = int(conf.MaxWsConnections)

    wsConf := conf.WsConf
    if wsConf == nil {
        wsConf = NewWsConf()
    }

    s.minByteRatePerSec = int(wsConf.WsMinByteRatePerSec)

    s.wsSH = &wsServerHandler{}
    s.wsSH.isBandWidthLimitSet = s.networkBandWidth > 0
    s.wsSH.wsConn = make(map[*wsConn]struct{})

    s.wsH = &wsHandler{}
    s.wsH.setDefault(wsConf)
    s.wsH.isMaskRequired = true
    s.wsH.wsSH = s.wsSH

    s.onMalformedRequest = conf.OnMalformedRequest
    s.onHttpRequest = conf.OnHttpRequest
    s.onWebsocketRequest = conf.OnWebsocketRequest

    return &s
}
