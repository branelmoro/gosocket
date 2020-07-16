package gosocket

import "net/url"

type HttpRequest interface {
	Method() string
	Host() string
	Uri() string
	Url() *url.URL
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
	headers map[string]string
	// conn *Conn

	url *url.URL

	// these fields are set if response is read
	bytes []byte
	headerStart int
	headerEnd int
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

func (r *httpRequest) Url() *url.URL {
	return r.url
}

func (r *httpRequest) Protocol() string {
	return r.protocol
}

func (r *httpRequest) Header() map[string]string {
	return r.headers
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
		r.protocol == "HTTP/1.1" &&
		// r.headers["origin"] != "" &&
		r.headers["upgrade"] != "" &&
		r.headers["upgrade"] == "websocket" &&
		r.headers["sec-websocket-version"] == "13" &&
		r.headers["sec-websocket-key"] != ""
}