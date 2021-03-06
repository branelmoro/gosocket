package gosocket

import (
	"net/url"
	"regexp"
	"strings"
)

func isWS(ch byte) bool {
	return ch == 0x20 || ch == 0x9
}

func isControlChar(ch byte) bool {
	return (ch < 0x20 && ch != 0x9) || ch == 0x7f
}

type httpReader struct {
	*Conn
	req *httpRequest
	readByteCount int
	headerLimit int
}

func (r *httpReader) readByte() (byte, error) {
	numBytes, readBytes, err := r.read(1)
	if err != nil {
		return 0, newReadError(err)
	}
	r.req.bytes = append(r.req.bytes, readBytes...)
	r.readByteCount += numBytes
	if r.readByteCount > 10000 {
		return 0, newReadError(err)
	}
	return readBytes[0], nil
}

func (r *httpReader) readRequest() (*httpRequest, error) {

	r.req = &httpRequest{}

	err := r.readInitialNewLines()
	if err != nil {
		return r.req, err
	}

	err = r.readRequestLine()
	if err != nil {
		return r.req, err
	}

	err = r.readHeader()
	return r.req, err
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

func (r *httpReader) readInitialNewLines() error {
	for {
		isNL, err := r.isLineBreak()
		if err != nil {
			return err
		}
		if !isNL {
			return nil
		}
	}
}

func (r *httpReader) readMethod() error {
	var(
		byteReceived byte
		isNL bool
		start int
		end int
		err error
	)
	start = len(r.req.bytes)-1
	end = len(r.req.bytes)
	for {
		isNL, err = r.isLineBreak()
		if err != nil {
			return err
		}
		if isNL {
			// unexpected end of request line
			return newHttpMalformedError("unexpected end of request line")
		}
		byteReceived = r.req.bytes[end]
		if byteReceived == 0x20 {
			// break at space character
			break
		}
		// valid characters
		end += 1
		if end-start == 7 {
			byteReceived, err = r.readByte()
			if err != nil {
				return err
			}
			if byteReceived != 0x20 {
				// return invalid http request error
				return newHttpMalformedError("invalid http request error")
			}
			break
		}
	}
	r.req.method = string(r.req.bytes[start:end])
	switch (r.req.method) {
		case "GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "TRACE", "CONNECT":
			break
		default:
			// invalid bytes received in http request method
			return newHttpMalformedError("invalid bytes received in http request method")
	}
	return err
}

func (r *httpReader) readRequestLine() error {
	var(
		prevByte byte
		isUriStarted bool
		isNL bool
		uriStart int
		uriEnd int
		protoStart int
		err error
	)
	err = r.readMethod()
	if err != nil {
		return err
	}
	isUriStarted = false
	uriStart = len(r.req.bytes)
	uriEnd = len(r.req.bytes)
	index := uriStart
	for {
		isNL, err = r.isLineBreak()
		if err != nil {
			return err
		}
		if isNL {
			// end of request line
			if !isUriStarted {
				// no uri found in request start line
				return newHttpMalformedError("no uri found in request start line")
			} else {
				protoStart = index-8
				if uriEnd >= protoStart {
					// error in request line string
					return newHttpMalformedError("error in request line string")
				}
			}
			break
		}
		byteReceived := r.req.bytes[index]
		if isUriStarted {
			if byteReceived == 0x20 {
				if prevByte != 0x20 {
					uriEnd = index
				}
			}
		} else {
			if byteReceived != 0x20 {
				isUriStarted = true
				uriStart = index
			}
		}
		prevByte = byteReceived
		index += 1
	}
	r.req.uri = string(r.req.bytes[uriStart:uriEnd])
	r.req.URL, err = url.Parse(r.req.uri)
	if err != nil {
		// return url parsing error
		return newHttpUriError(err)
	}
	if r.req.URL.Host != "" {
		r.req.host = r.req.URL.Host
		if !r.isValidHttpHost(r.req.host) {
			// return invalid host received in http request
			return newHttpMalformedError("invalid host received in http request")
		}
	}
	r.req.protocol = string(r.req.bytes[protoStart:index])
	if r.req.protocol != "HTTP/1.1" {
		// return invalid http protocol error, only HTTP/1.1 allowed
		return newHttpMalformedError("invalid http protocol error, only HTTP/1.1 allowed")
	}
	return err
}

func (r *httpReader) readHeader() error {
	var(
		isNL bool
		field string
		fieldStart int
		fieldEnd int
		valStart int
		valEnd int
		err error
	)
	r.req.headerStart = len(r.req.bytes)
	r.headerLimit = r.req.headerStart + r.server.httpMaxHeaderSize
	field = ""
	r.req.header = make(map[string]string)
	for {
		isNL, err = r.isLineBreak()
		if err != nil {
			return err
		}
		if isNL {
			// finished reading header
			break
		}
		if isWS(r.req.bytes[len(r.req.bytes) - 1]) {
			// header line contains folded value
			if field == "" {
				// whitespace found at start of header field
				return newHttpMalformedError("whitespace found at start of header field")
			} else {
				err = r.readHeaderBytes(false)
				if err != nil {
					return err
				}
				valEnd = len(r.req.bytes)
			}
		} else {
			if field != "" {
				r.setHeaderField(field, valStart, valEnd)
				if field == "host" {
					if r.req.host != "" {
						// return host received in both uri and http header
						return newHttpMalformedError("host received in both uri and http header")
					}
					if !r.isValidHttpHost(r.req.header[field]) {
						// return invalid host received in http
						return newHttpMalformedError("invalid host received in http")
					}
					r.req.host = r.req.header[field]
				}
			}
			fieldStart = len(r.req.bytes) - 1
			err = r.readHeaderBytes(true)
			if err != nil {
				return err
			}
			fieldEnd = len(r.req.bytes) - 1
			field = strings.TrimSpace(strings.ToLower(string(r.req.bytes[fieldStart:fieldEnd])))
			valStart = len(r.req.bytes)
			valEnd = valStart
		}
	}
	if field != "" {
		r.setHeaderField(field, valStart, valEnd)
		if field == "host" {
			if r.req.host != "" {
				// return invalid host received in both uri and http header
				return newHttpMalformedError("invalid host received in both uri and http header")
			}
			if !r.isValidHttpHost(r.req.header[field]) {
				// return invalid host received in http
				return newHttpMalformedError("invalid host received in http")
			}
			r.req.host = r.req.header[field]
		}
	}
	r.req.headerEnd = len(r.req.bytes) - 2

	if r.req.host == "" {
		// no host found in uri or header
		return newHttpMalformedError("no host found in uri or header")
	}

	return err
}


func (r *httpReader) setHeaderField(field string, start int, end int) {
	space := regexp.MustCompile(`\s+`)
	if _, isPresent := r.req.header[field]; isPresent {
		r.req.header[field] += ", " + space.ReplaceAllString(strings.TrimSpace(string(r.req.bytes[start:end])), " ")
	} else {
		r.req.header[field] = space.ReplaceAllString(strings.TrimSpace(string(r.req.bytes[start:end])), " ")
	}
}

func (r *httpReader) isValidHttpHost(val string) bool {
	data := strings.SplitN(val, ":", 2);
	return data[0] == r.server.host
}

func (r *httpReader) readHeaderBytes(isField bool) error {
	var(
		byteReceived byte
		isNL bool
		index int
		err error
	)
	index = len(r.req.bytes)
	for {
		isNL, err = r.isLineBreak()
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
			byteReceived = r.req.bytes[index]
			if isField && byteReceived == 0x3a {
				// stop at ":", end of header field
				return nil
			}
			index += 1
			if index > r.headerLimit {
				// header size is more than allowed size
				return newHttpMalformedError("header size is more than allowed size")
			}
		}
	}
	return err
}
