package gosocket

import(
    "fmt"
    "github.com/mailru/easygo/netpoll"
    "sync"
    // "sync/atomic"
    "time"
)

type wsServerHandler struct {

    _rOpsLock sync.Mutex
    cntReadOps uint
    _wOpsLock sync.Mutex
    cntWriteOps uint

    wsConn map[*wsConn]struct{}

    isBandWidthLimitSet bool

    poller netpoll.Poller
}

func (wc *wsServerHandler) addConn(ws *wsConn) {
    wc.wsConn[ws] = struct{}{}
}

func (wc *wsServerHandler) delConn(ws *wsConn) {
    if _, ok := wc.wsConn[ws]; ok {
        delete(wc.wsConn, ws)
    }
}

func (wc *wsServerHandler) wsCount() int {
    return len(wc.wsConn)
}

func (wc *wsServerHandler) addReadOps() {
    wc._rOpsLock.Lock()
    defer wc._rOpsLock.Unlock()
    wc.cntReadOps++
    // AddUintptr(&wc.cntReadOps, delta uintptr) (new uintptr)
}

func (wc *wsServerHandler) delReadOps() {
    wc._rOpsLock.Lock()
    defer wc._rOpsLock.Unlock()
    wc.cntReadOps--
}

func (wc *wsServerHandler) addWriteOps() {
    wc._wOpsLock.Lock()
    defer wc._wOpsLock.Unlock()
    wc.cntWriteOps++
}

func (wc *wsServerHandler) delWriteOps() {
    wc._wOpsLock.Lock()
    defer wc._wOpsLock.Unlock()
    wc.cntWriteOps--
}

func (wc *wsServerHandler) state() {

}

func (wc *wsServerHandler) shutDown() {
    msg := NewCloseMsg(CC_GOING_AWAY, "server shutting down")
    for wc.wsCount() > 0 {
        var(
            ws *wsConn
        )
        parallelClose := 1000
        counter := 0
        for k, _ := range wc.wsConn {
            counter++
            if counter%parallelClose == 1 {
                ws = k
            } else {
                go k.writer().Close(msg)
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
        fmt.Println("Please wait...", wc.wsCount(), " connections are still open")
        time.Sleep(200 * time.Millisecond)
    }
    fmt.Println("All connections closed successfully")
}
