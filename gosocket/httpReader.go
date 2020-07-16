package gosocket

import (
	// "io"
	"net"
	// "net/url"
	"regexp"
	"strings"
	"time"
)

type httpConn interface {
	close() error
}

type httpReaderConn interface {
	httpConn
	setReadTimeOut(time.Time) error
	read([]byte) (int, error)
}

type httpWriterConn interface {
	httpConn
	setWriteTimeOut(time.Time) error
	write([]byte) (int, error)
}


type httpconn struct {
	net.Conn
}

func (c *httpconn) setReadTimeOut(t time.Time) error {
	return c.SetReadDeadline(t)
}

func (c *httpconn) setWriteTimeOut(t time.Time) error {
	return c.SetWriteDeadline(t)
}

func (c *httpconn) read(data []byte) (int, error) {
	return c.Read(data)
}

func (c *httpconn) write(data []byte) (int, error) {
	return c.Write(data)
}

func (c *httpconn) close() error {
	return c.Close()
}


type httpReader struct {
	conn httpReaderConn
	readBytes []byte
	maxHeaderSize int
}

func (r *httpReader) read(size int) (int, []byte, error) {
	buffer := make([]byte, size)
	num_bytes, err := r.conn.read(buffer)
	return num_bytes, buffer[:num_bytes], err
}

func (r *httpReader) readByte() (byte, error) {
	_, readBytes, err := r.read(1)
	if err != nil {
		return 0, newReadError(err)
	}
	r.readBytes = append(r.readBytes, readBytes...)

	if len(r.readBytes) > r.maxHeaderSize {
		// header size is more than allowed size
		return 0, newHttpMalformedError("header size is more than allowed size")
	}
	return readBytes[0], nil
}

func (r *httpReader) isLineBreak() (bool, error) {
	byteReceived, err := r.readByte()
	if err != nil {
		return false, err
	}
	if byteReceived == 0x0d {
		byteReceived, err = r.readByte()
		if err != nil {
			return false, err
		}
		if byteReceived == 0x0a {
			return true, nil
		} else {
			// invalid line break, only \r control character received
			return false, newHttpMalformedError("invalid line break, only \r control character received")
		}
	}
	if isControlChar(byteReceived) {
		// control char received in request
		return false, newHttpMalformedError("control char received in request")
	}
	return false, nil
}

func (r *httpReader) readHeader(headerSetter func(string, string) error) error {
	var(
		// isNL bool
		// field string
		// fieldStart int
		// fieldEnd int
		valStart int
		valEnd int
		err error
	)
	field := ""
	for {
		isNL, err := r.isLineBreak()
		if err != nil {
			return err
		}
		if isNL {
			// finished reading header
			break
		}
		if isWhiteSpace(r.readBytes[len(r.readBytes) - 1]) {
			// header line contains folded value
			if field == "" {
				// whitespace found at start of header field
				return newHttpMalformedError("whitespace found at start of header field")
			} else {
				err = r.readHeaderBytes(false)
				if err != nil {
					return err
				}
				valEnd = len(r.readBytes)
			}
		} else {
			if field != "" {
				val := r.getTrimmedHeaderVal(r.readBytes[valStart:valEnd])
				err = headerSetter(field, val)
				if err != nil {
					return err
				}
			}
			fieldStart := len(r.readBytes) - 1
			err = r.readHeaderBytes(true)
			if err != nil {
				return err
			}
			fieldEnd := len(r.readBytes) - 1
			field = strings.TrimSpace(strings.ToLower(string(r.readBytes[fieldStart:fieldEnd])))
			valStart = len(r.readBytes)
			valEnd = valStart
		}
	}
	if field != "" {
		val := r.getTrimmedHeaderVal(r.readBytes[valStart:valEnd])
		err = headerSetter(field, val)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *httpReader) getTrimmedHeaderVal(val []byte) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(strings.TrimSpace(string(val)), " ")
}

func (r *httpReader) readHeaderBytes(isField bool) error {
	var(
		// byteReceived byte
		// isNL bool
		// index int
		// err error
	)
	index := len(r.readBytes)
	for {
		isNL, err := r.isLineBreak()
		if err != nil {
			return err
		}
		if isNL {
			if isField {
				// line break found in header field
				return newHttpMalformedError("line break found in header field")
			} else {
				return nil
			}
		} else {
			byteReceived := r.readBytes[index]
			if isField && byteReceived == 0x3a {
				// stop at ":", end of header field
				return nil
			}
			index += 1
		}
	}
	return nil
}
