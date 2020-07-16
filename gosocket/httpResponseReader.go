package gosocket

import (
	"strconv"
)

type httpResponseReader struct {
	*httpReader
}

func (r *httpResponseReader) readResponse() (*httpResponse, error) {

	res := &httpResponse{}

	err := r.readResponseLine(res)
	if err != nil {
		res.bytes = r.readBytes
		return res, err
	}

	res.headerStart = len(r.readBytes)

	headerSetter := func(field string, val string) error {
		if _, isPresent := res.headers[field]; isPresent {
			res.headers[field] += ", " + val
		} else {
			res.headers[field] = val
		}
		return nil
	}

	err = r.readHeader(headerSetter)
	res.bytes = r.readBytes
	if err != nil {
		return res, err
	}
	res.headerEnd = len(r.readBytes) - 2

	return res, err
}

func (r *httpResponseReader) readHttpVersion(res *httpResponse) error {
	var(
		// byteReceived byte
		// isNL bool
		// end int
		// err error
	)
	end := 0
	for {
		isNL, err := r.isLineBreak()
		if err != nil {
			return err
		}
		if isNL {
			// unexpected end of response line
			return newHttpMalformedError("unexpected end of response line")
		}
		byteReceived := r.readBytes[end]
		if byteReceived == 0x20 {
			// break at space character
			break
		}
		// valid characters
		end += 1
		if end == 8 {
			byteReceived, err = r.readByte()
			if err != nil {
				return err
			}
			if byteReceived != 0x20 {
				// return invalid http response error
				return newHttpMalformedError("invalid http response error")
			}
			break
		}
	}
	res.protocol = string(r.readBytes[:end])

	if res.protocol != "HTTP/1.1" {
		// return invalid http protocol error, only HTTP/1.1 allowed
		return newHttpMalformedError("invalid http protocol error, only HTTP/1.1 allowed")
	}
	return nil
}

func (r *httpResponseReader) readStatusCode(res *httpResponse) error {
	var(
		// byteReceived byte
		// isNL bool
		// end int
		err error
	)
	start := len(r.readBytes)
	end := len(r.readBytes)
	for {
		isNL, err := r.isLineBreak()
		if err != nil {
			return err
		}
		if isNL {
			// unexpected end of response line
			return newHttpMalformedError("unexpected end of response line")
		}
		byteReceived := r.readBytes[end]
		if byteReceived == 0x20 {
			// break at space character
			break
		}
		// valid characters
		end += 1
		if end == 3 {
			byteReceived, err = r.readByte()
			if err != nil {
				return err
			}
			if byteReceived != 0x20 {
				// return invalid http response error
				return newHttpMalformedError("invalid http response error")
			}
			break
		}
	}

	res.code, err = strconv.Atoi(string(r.readBytes[start:end]))

	return err
}

func (r *httpResponseReader) readResponseLine(res *httpResponse) error {
	err := r.readHttpVersion(res)
	if err != nil {
		return err
	}
	err = r.readStatusCode(res)
	if err != nil {
		return err
	}

	// read till line break
	start := len(r.readBytes)
	end := len(r.readBytes)
	for {
		isNL, err := r.isLineBreak()
		if isNL {
			break
		}
		if err != nil {
			return err
		}
		end += 1
	}

	res.reasonPhrase = string(r.readBytes[start:end])
	return nil
}
