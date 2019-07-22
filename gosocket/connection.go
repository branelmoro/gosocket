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
}

func (c *Conn) fdesc() *netpoll.Desc {
	return c.desc
}

func (c *Conn) fpoller() netpoll.Poller {
	return c.poller
}

func (c *Conn) setIoTimeOut(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *Conn) setReadTimeOut(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Conn) setWriteTimeOut(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func (c *Conn) read(size int) (int, []byte, error) {
	var(
		err error
		num_bytes int
	)
	buffer := make([]byte, size)
	num_bytes, err = c.conn.Read(buffer)
	return num_bytes, buffer[:num_bytes], err
}

func (c *Conn) write(data []byte) (int, error) {
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
