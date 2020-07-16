package gosocket

import (
	"net/url"
	"strings"
)

type httpRequestReader struct {
	*httpReader
	host string
}

func (r *httpRequestReader) readRequest() (*httpRequest, error) {

	req := &httpRequest{}

	err := r.readInitialNewLines()
	if err != nil {
		req.bytes = r.readBytes
		return req, err
	}

	err = r.readRequestLine(req)
	if err != nil {
		req.bytes = r.readBytes
		return req, err
	}

	req.headerStart = len(r.readBytes)

	headerSetter := func(field string, val string) error {
		if _, isPresent := req.headers[field]; isPresent {
			req.headers[field] += ", " + val
		} else {
			req.headers[field] = val
		}
		if field == "host" {
			if req.host != "" {
				// return invalid host received in both uri and http headers
				return newHttpMalformedError("host received in both uri and http headers")
			}
			if !r.isValidHttpHost(req.headers[field]) {
				// return invalid host received in http
				return newHttpMalformedError("invalid host received in http")
			}
		}
		return nil
	}

	err = r.readHeader(headerSetter)
	req.bytes = r.readBytes
	if err != nil {
		return req, err
	}
	req.headerEnd = len(r.readBytes) - 2

	if req.host == "" {
		if _, isPresent := req.headers["host"]; !isPresent {
			// no host found in uri or headers
			return req, newHttpMalformedError("no host found in uri or headers")
		} else {
			req.host = req.headers["host"]
		}
	}

	return req, err
}

func (r *httpRequestReader) readInitialNewLines() error {
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

func (r *httpRequestReader) readMethod(req *httpRequest) error {
	var(
		// byteReceived byte
		// isNL bool
		// start int
		// end int
		// err error
	)
	start := len(r.readBytes)-1
	end := len(r.readBytes)
	for {
		isNL, err := r.isLineBreak()
		if err != nil {
			return err
		}
		if isNL {
			// unexpected end of request line
			return newHttpMalformedError("unexpected end of request line")
		}
		byteReceived := r.readBytes[end]
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
	req.method = string(r.readBytes[start:end])
	switch (req.method) {
		case "GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "TRACE", "CONNECT":
			break
		default:
			// invalid bytes received in http request method
			return newHttpMalformedError("invalid bytes received in http request method")
	}
	return nil
}

func (r *httpRequestReader) readRequestLine(req *httpRequest) error {
	var(
		prevByte byte
		// isUriStarted bool
		// isNL bool
		// uriStart int
		// uriEnd int
		protoStart int
		// err error
	)
	err := r.readMethod(req)
	if err != nil {
		return err
	}
	isUriStarted := false
	uriStart := len(r.readBytes)
	uriEnd := len(r.readBytes)
	index := uriStart
	for {
		isNL, err := r.isLineBreak()
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
		byteReceived := r.readBytes[index]
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
	req.uri = string(r.readBytes[uriStart:uriEnd])
	req.url, err = url.Parse(req.uri)
	if err != nil {
		// return url parsing error
		return newHttpUriError(err)
	}
	if req.url.Host != "" {
		req.host = req.url.Host
		if !r.isValidHttpHost(req.host) {
			// return invalid host received in http request
			return newHttpMalformedError("invalid host received in http request")
		}
	}
	req.protocol = string(r.readBytes[protoStart:index])
	if req.protocol != "HTTP/1.1" {
		// return invalid http protocol error, only HTTP/1.1 allowed
		return newHttpMalformedError("invalid http protocol error, only HTTP/1.1 allowed")
	}
	return nil
}

func (r *httpRequestReader) isValidHttpHost(val string) bool {
	data := strings.SplitN(val, ":", 2);
	return data[0] == r.host
}