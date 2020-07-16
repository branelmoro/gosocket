package gosocket

import (
    "net"
    "time"
    "github.com/mailru/easygo/netpoll"
)


type netpollConn interface {
	setReadTimeOut(time.Time) error
	setWriteTimeOut(time.Time) error
	read(int) (int, []byte, error)
	write([]byte) (int, error)
	close() error

	stopPoller() error

	wsH() *wsHandler
}

type conn struct {
    c    net.Conn
    desc    *netpoll.Desc
    poller  netpoll.Poller
    handler *wsHandler
}

func (c *conn) wsH() *wsHandler {
	return c.handler
}

// func (c *conn) fdesc() *netpoll.Desc {
// 	return c.desc
// }

// func (c *conn) fpoller() netpoll.Poller {
// 	return c.poller
// }

func (c *conn) setReadTimeOut(t time.Time) error {
	return c.c.SetReadDeadline(t)
}

func (c *conn) setWriteTimeOut(t time.Time) error {
	return c.c.SetWriteDeadline(t)
}

func (c *conn) readBytes(size int) (int, []byte, error) {
	buffer := make([]byte, size)
	num_bytes, err := c.c.Read(buffer)
	return num_bytes, buffer[:num_bytes], err
}

func (c *conn) close() error {
	c.stopPoller()
	return c.c.Close()
}

func (c *conn) stopPoller() error {
	p := c.poller
	p.Stop(c.desc)
	return c.desc.Close()
}



type serverConn struct {
	*conn
}

func (sc *serverConn) controlRead(size int) (int, []byte, error) {
	var(
		err error
		totolCount int
		readCount int
		data []byte
		readData []byte
		startTime time.Time
	)

	totolCount = 0
	for size > 0 {
		startTime = time.Now()
		if size > sc.handler.maxByteRate {
			readCount = sc.handler.maxByteRate
		} else {
			readCount = size
		}
		readCount, data, err = sc.readBytes(sc.handler.maxByteRate)
		totolCount += readCount
		readData = append(readData, data...)
		if err != nil {
			return totolCount, readData, err
		}
		size -= readCount
		sc.controlSpeed(readCount, startTime)
	}

	return totolCount, readData, err
}

func (sc *serverConn) controlSpeed(readCount int, startTime time.Time) {
	maxExpectedCount := int((float64(time.Now().Sub(startTime).Nanoseconds())/float64(time.Second)) * float64(sc.handler.maxByteRate))
	if readCount > maxExpectedCount {
		time.Sleep(time.Duration(((readCount-maxExpectedCount)/sc.handler.maxByteRate) * int(time.Second)))
	}
}

func (sc *serverConn) read(size int) (int, []byte, error) {
	sc.handler.wsSH.addReadOps()
	defer sc.handler.wsSH.delReadOps()
	if sc.handler.wsSH.isBandWidthLimitSet {
		return sc.controlRead(size)
	} else {
		return sc.readBytes(size)
	}
}

func (sc *serverConn) write(data []byte) (int, error) {
	sc.handler.wsSH.addWriteOps()
	defer sc.handler.wsSH.delWriteOps()
	return sc.c.Write(data)
}



type clientConn struct {
	*conn
}

func (cc *clientConn) read(size int) (int, []byte, error) {
	return cc.readBytes(size)
}

func (cc *clientConn) write(data []byte) (int, error) {
	return cc.c.Write(data)
}