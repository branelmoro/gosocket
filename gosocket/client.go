package gosocket

import(
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/mailru/easygo/netpoll"
)

type client struct {
    poller netpoll.Poller
    wsH *wsHandler
    headers map[string]string
}

func (cl *client) getWebSocketHeaders() map[string]string {
	headers := make(map[string]string)
	for key, val := range cl.headers {
		headers[strings.ToLower(key)] = val
	}
	if _, ok := headers["upgrade"]; !ok {
		headers["upgrade"] = "websocket"
	}
	if _, ok := headers["sec-websocket-version"]; !ok {
		headers["sec-websocket-version"] = "13"
	}
	if _, ok := headers["sec-websocket-key"]; !ok {
		headers["sec-websocket-key"] = "yhagsdhkashdkjnask"
	}
	return headers
}

func (cl *client) sendRequest(c net.Conn, data []byte) error {
	startIndex := 0
	for startIndex != len(data) {
		numBytes, err := c.Write(data[startIndex:])
		if err != nil {
			return err
		}
		startIndex += numBytes
	}
	return nil
}

func (cl *client) connectToHost(host string, port int, isSecured bool) (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
}

func (cl *client) Connect(uri string) (WsWriter, error) {
	url, err := url.Parse(uri)
	if err != nil {
		// return url parsing error
		return nil, fmt.Errorf("Invalid Connection URI")
	}

	host := url.Hostname()
	if host == "" {
		return nil, fmt.Errorf("Invalid host in Connection URI")
	}

	// set default port for http
	port := 80
	isSecured := false

	scheme := url.Scheme

	if scheme == "ws" {
		url.Scheme = "http"
	}

	if scheme == "wss" {
		url.Scheme = "https"
	}

	switch url.Scheme {
		case "http":
			break
		case "https":
			port = 443
			isSecured = true
			break
		default:
			return nil, fmt.Errorf("Invalid Protocol - %s in URI", scheme)
	}

	portStr := url.Port()
	if portStr != "" {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("Invalid Port - %s in URI", portStr)
		}
	}

	c, err := cl.connectToHost(host, port, isSecured)
	if err != nil {
		return nil, err
	}


	// create websocket request
	wsHeaders := cl.getWebSocketHeaders()
	wsReqBytes := []byte(fmt.Sprintf("GET %s HTTP/1.1\r\n", uri))
	for key, val := range wsHeaders {
		wsReqBytes = append(wsReqBytes, []byte(strings.Title(key) + ": " + val + "\r\n")...)
	}
	wsReqBytes =  append(wsReqBytes, 0xd, 0xa)

	err = cl.sendRequest(c, wsReqBytes)
	if err != nil {
		return nil, err
	}


    hConn := &httpconn{c}
    reader := &httpResponseReader{&httpReader{ conn: hConn }}
    res, err := reader.readResponse()
    if err != nil {
    	return nil, err
    }

    if res.code != 101 {
		return nil, fmt.Errorf("Invalid response code - %d received, Websocket failed!", res.code)
    }

	if res.code == 101 &&
	res.protocol == "HTTP/1.1" &&
	// r.headers["origin"] != "" &&
	// r.headers["upgrade"] != "" &&
	res.headers["connection"] == "Upgrade" &&
	res.headers["upgrade"] == "websocket" &&
	res.headers["sec-websocket-version"] == "13" &&
	res.headers["sec-websocket-accept"] != "" &&
	res.headers["sec-websocket-accept"] == generateSecWebSocketAccept(wsHeaders["sec-websocket-key"]) {

		// evts := netpoll.EventOneShot | netpoll.EventPollerClosed | netpoll.EventErr | netpoll.EventWriteHup | netpoll.EventReadHup | netpoll.EventHup | netpoll.EventRead | netpoll.EventWrite | netpoll.EventEdgeTriggered
		evts := netpoll.EventPollerClosed | netpoll.EventWriteHup | netpoll.EventReadHup | netpoll.EventHup | netpoll.EventRead | netpoll.EventEdgeTriggered

		// Get netpoll descriptor with EventRead|EventEdgeTriggered.
		desc := netpoll.Must(netpoll.Handle(c, evts))

		poller := cl.poller

		wc := wsConn{
			netpollConn: &clientConn{ &conn{ c, desc, poller, cl.wsH } },
			// ConnData: options.WsData,
			// flate: flate,
		}
		openWebSocket(&wc, poller, desc)
		return wc.writer(), nil
	} else {
		return nil, fmt.Errorf("Connection upgrade to websocket failed!")
	}

}

type Client interface {
    Connect(string) (WsWriter, error)
}


func NewClient(conf *ClientConf) Client {
	var err error

    cl := &client{}

    cl.poller, err = netpoll.New(nil)
    if err != nil {
		panic("Error in creating connection poller")
    }
    if conf == nil {
        conf = NewClientConf()
    }
    wsConf := conf.WsConf
    if wsConf == nil {
        wsConf = NewWsConf()
    }

    cl.wsH = &wsHandler{}
    cl.wsH.setDefault(wsConf)
    cl.wsH.isMaskRequired = false


    cl.headers = conf.Headers
    return cl
}
