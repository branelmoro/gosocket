package gosocket

// WsWriter describes connection object that implements all functions
// for sending data over websocket connection
type WsWriter interface {

	// Send sends message over websocket connetcion.
	// it returns error object is anything went wrong in sending data.
	Send(Message) error

	// Send sends text data over websocket connetcion.
	// it returns error object is anything went wrong in sending data.
	SendText(string) error

	// Send sends binary data over websocket connetcion.
	// it returns error object is anything went wrong in sending data.
	SendBinary([]byte) error

	// Closes websocket Connetion.
	// it returns error object is anything went wrong in closing connection.
	Close(CloseMsg) error
}



type Server interface {

    Run()

    Shutdown()

    Restart()

    // DisconnectAll() error

    // StopAccept() error

    // StartAccept() error

    // StopRead() error

    // StartRead() error

    // StopWrite() error

    // StartWrite() error
}

// things to include:
// bandwidth
// maxConnections
// allowed_origins
// connection pool


// no of servers - n
// no of clustor connections per server = n-1
// total no of clients - c
// clients connections per server = sc = c/n
// incoming data on server connection = sc/n
// outgoing traffic on cluster connection = cc = sc - sc/n

// cc/sc = (n-1)/n

// cc + sc = bw

// cc = sc - sc/n

// bw - sc = sc - sc/n

// bw = sc(2-1/n)

// sc = bw/(2-1/n)

// sc = (bw * n)/(2n-1)


