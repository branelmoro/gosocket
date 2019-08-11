package gosocket

import (
    "net"
    "time"
    "github.com/mailru/easygo/netpoll"
)

// Conn represents single connection instance.
type Conn struct {
    conn    net.Conn
    desc    *netpoll.Desc
    poller  netpoll.Poller
    server *server
    speedControl bool
}

func (c *Conn) fdesc() *netpoll.Desc {
	return c.desc
}

func (c *Conn) fpoller() netpoll.Poller {
	return c.poller
}

func (c *Conn) setReadTimeOut(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Conn) setWriteTimeOut(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func (c *Conn) controlRead(size int) (int, []byte, error) {
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
		if size > c.server.maxIOSpeed() {
			readCount = c.server.maxIOSpeed()
		} else {
			readCount = size
		}
		readCount, data, err = c.readBytes(c.server.maxIOSpeed())
		totolCount += readCount
		readData = append(readData, data...)
		if err != nil {
			return totolCount, readData, err
		}
		size -= readCount
		c.controlSpeed(readCount, startTime)
	}

	return totolCount, readData, err
}

func (c *Conn) controlSpeed(readCount int, startTime time.Time) {
	maxExpectedCount := int((float64(time.Now().Sub(startTime).Nanoseconds())/float64(time.Second)) * float64(c.server.maxIOSpeed()))
	if readCount > maxExpectedCount {
		time.Sleep(time.Duration(((readCount-maxExpectedCount)/c.server.maxIOSpeed()) * int(time.Second)))
	}
}

func (c *Conn) readBytes(size int) (int, []byte, error) {
	buffer := make([]byte, size)
	num_bytes, err := c.conn.Read(buffer)
	return num_bytes, buffer[:num_bytes], err
}

func (c *Conn) read(size int) (int, []byte, error) {
	c.server.addReadOps()
	defer c.server.delReadOps()
	if c.speedControl {
		return c.controlRead(size)
	} else {
		return c.readBytes(size)
	}
}

func (c *Conn) write(data []byte) (int, error) {
	c.server.addWriteOps()
	defer c.server.delWriteOps()
	var(
		err error
		num_bytes int
	)
	num_bytes, err = c.conn.Write(data)
	return num_bytes, err
}

func (c *Conn) close() error {
	p := c.poller
	p.Stop(c.desc)
	c.desc.Close()
	err := c.conn.Close()
	return err
}
