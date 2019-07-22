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
