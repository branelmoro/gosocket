package gosocket

import (
	"fmt"
	"strings"
)

type HttpResponse interface {
	Protocol() string
	Code() int
	Header() map[string]string
	RawHeader() []byte
	Raw() []byte
	ToString() string
}

type httpResponse struct {
	protocol string
	code int
	reasonPhrase string
	headers map[string]string

	// these fields are set if response is read
	bytes []byte
	headerStart int
	headerEnd int
}

func (r *httpResponse) Protocol() string {
	return r.protocol
}

func (r *httpResponse) Code() int {
	return r.code
}

func (r *httpResponse) ReasonPhrase() string {
	return r.reasonPhrase
}

func (r *httpResponse) Header() map[string]string {
	return r.headers
}

func (r *httpResponse) RawHeader() []byte {
	return r.bytes[r.headerStart:r.headerEnd]
}

func (r *httpResponse) Raw() []byte {
	return r.bytes
}

func (r *httpResponse) ToString() string {
	return string(r.bytes)
}

func (r *httpResponse) toBytes() []byte {
	bytes := []byte(fmt.Sprintf(r.protocol + " %d %s\r\n", r.code, httpStatusText[r.code]))
	for key, val := range r.headers {
		bytes = append(bytes, []byte(strings.Title(key) + ": " + val + "\r\n")...)
	}
	return append(bytes, 0xd, 0xa)
}
