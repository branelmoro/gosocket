package gosocket

import "net/url"

type HttpRequest interface {
	Method() string
	Host() string
	Uri() string
	Protocol() string
	Header() map[string]string
	RawHeader() []byte
	Raw() []byte
	ToString() string
	isWebSocketRequest() bool
}

type httpRequest struct {
	method string
	host string
	uri string
	protocol string
	header map[string]string
	conn *Conn

	bytes []byte
	headerStart int
	headerEnd int

	URL *url.URL
}

func (r *httpRequest) Method() string {
	return r.method
}

func (r *httpRequest) Host() string {
	return r.method
}

func (r *httpRequest) Uri() string {
	return r.uri
}

func (r *httpRequest) Protocol() string {
	return r.protocol
}

func (r *httpRequest) Header() map[string]string {
	return r.header
}

func (r *httpRequest) RawHeader() []byte {
	return r.bytes[r.headerStart:r.headerEnd]
}

func (r *httpRequest) Raw() []byte {
	return r.bytes
}

func (r *httpRequest) ToString() string {
	return string(r.bytes)
}

func (r *httpRequest) isWebSocketRequest() bool {
	return r.method == "GET" &&
		// r.header["origin"] != "" &&
		r.header["upgrade"] != "" &&
		r.header["upgrade"] == "websocket" &&
		r.header["sec-websocket-version"] == "13" &&
		r.header["sec-websocket-key"] != ""
}