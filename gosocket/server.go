package gosocket

import (
    "crypto/rand"
    "fmt"
    "crypto/tls"
    "net"
    "sync"
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

    isRunning bool

    httpRquestTimeOut time.Duration
    httpMaxRequestLineSize int
    httpMaxHeaderSize int

    wsMaxFrameSize int
    wsMaxMessageSize int

    wsHeaderReadTimeout time.Duration
    wsMinByteRatePerSec int
    wsCloseTimeout time.Duration

    networkBandWidth int
    maxWsConnection int

    wsConn map[*wsConn]interface{}

    cntHttpWrite uint
    cntWsWrite uint

    cntHttpRead uint
    cntWsRead uint

    _rOpsLock sync.Mutex
    cntReadOps uint
    _wOpsLock sync.Mutex
    cntWriteOps uint

    ioSpeed int
}

func (s *server) addConn(ws *wsConn) {
    s.wsConn[ws] = nil
}

func (s *server) delConn(ws *wsConn) {
    if _, ok := s.wsConn[ws]; ok {
        delete(s.wsConn, ws)
    }
}

func (s *server) wsCount() int {
    return len(s.wsConn)
}

func (s *server) addReadOps() {
    s._rOpsLock.Lock()
    defer s._rOpsLock.Unlock()
    s.cntReadOps++
}

func (s *server) delReadOps() {
    s._rOpsLock.Lock()
    defer s._rOpsLock.Unlock()
    s.cntReadOps--
}

func (s *server) addWriteOps() {
    s._wOpsLock.Lock()
    defer s._wOpsLock.Unlock()
    s.cntWriteOps++
}

func (s *server) delWriteOps() {
    s._wOpsLock.Lock()
    defer s._wOpsLock.Unlock()
    s.cntWriteOps--
}

func (s *server) maxIOSpeed() int {
    return s.ioSpeed
}

func (s *server) forever() {
    if s.networkBandWidth <= 0 {
        return
    }
    bw := s.networkBandWidth * 1024 * 1024
    for {
        iOps := s.cntReadOps + s.cntWriteOps
        if bw == 0 {
            s.ioSpeed = bw
        } else {
            s.ioSpeed = bw/int(iOps)
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

func (s *server) disconnectAll(msg CloseMsg) {
    var(
        ws *wsConn
    )
    parallelClose := 1000
    counter := 0
    for k, _ := range s.wsConn {
        counter++
        if counter%parallelClose == 1 {
            ws = k
        } else {
            go func(){
                k.writer().Close(msg)
            }()
        }
        if counter%parallelClose == 0 {
            ws.writer().Close(msg)
            ws = nil
            fmt.Println(counter, " connection closed")
        }
    }
    if ws != nil {
        ws.writer().Close(msg)
        ws = nil
        fmt.Println(counter, " connection closed")
    }

    for s.wsCount() > 0 {
        fmt.Println("Please wait...", s.wsCount(), " connections are still open")
        time.Sleep(200 * time.Millisecond)
    }
    fmt.Println("All connections closed successfully")
}

func (s *server) Run() {
    var(
        err error
        conn net.Conn
        poller netpoll.Poller
    )
    err = s.startListener()
    if err != nil {
        s.logError(err)
        fmt.Println("Error in starting listener: ", err)
        return
    }

    go s.forever()

    speedControl := s.networkBandWidth > 0

    // mark server running
    s.isRunning = true

    for {
        // wait if connection limit is reached
        for s.isRunning && s.maxWsConnection > 0 && s.wsCount() >= s.maxWsConnection {
            fmt.Println("Max connection limit reached, waiting for one second--- limit", s.maxWsConnection)
            time.Sleep(time.Second)
        }

        // Listen for an incoming connection.
        conn, err = s.listener.Accept()
        if err != nil {
            s.logError(err)
            break
        }

        poller, err = netpoll.New(nil)
        if err != nil {
            s.logError(err)
            conn.Close()
            continue
        }

        // Get netpoll descriptor with EventRead|EventEdgeTriggered.
        desc := netpoll.Must(netpoll.Handle(conn, netpoll.EventRead | netpoll.EventEdgeTriggered))

        socketConn := Conn{conn: conn, desc: desc, poller: poller, server: s, speedControl: speedControl}

        go s.handleConnection(&socketConn)
    }

    // close all opened connections
    msg := NewCloseMsg(CC_GOING_AWAY, "server shutting down")
    s.disconnectAll(msg)

}

func (s *server) handleConnection(conn *Conn) {
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

func (s *server) Shutdown() {
    err := s.listener.Close()
    if err != nil {
        s.logError(err)
        fmt.Println("Error Server Shutdown: ", err)
    }

    s.isRunning = false
}

func (s *server) Restart() {
    s.Shutdown()
    s.Run()
}


func NewServer(conf *ServerConf) Server {

    if conf == nil {
        conf = NewConf()
    }

    s := server{}

    s.host = conf.Host
    s.bindHosts = conf.BindHosts
    s.port = int(conf.Port)

    s.certPrivate = conf.CertPrivate
    s.certPublic = conf.CertPublic

    s.httpRquestTimeOut = conf.HttpRquestTimeOut
    s.httpMaxRequestLineSize = int(conf.HttpMaxRequestLineSize)
    s.httpMaxHeaderSize = int(conf.HttpMaxHeaderSize)

    s.wsMaxFrameSize = int(conf.WsMaxFrameSize)
    s.wsMaxMessageSize = int(conf.WsMaxMessageSize)

    s.wsHeaderReadTimeout = conf.WsHeaderReadTimeout
    s.wsMinByteRatePerSec = int(conf.WsMinByteRatePerSec)
    s.wsCloseTimeout = conf.WsCloseTimeout

    s.networkBandWidth = int(conf.NetworkBandWidth)
    s.maxWsConnection = int(conf.MaxWsConnection)

    s.wsConn = make(map[*wsConn]interface{})

    return &s
}
