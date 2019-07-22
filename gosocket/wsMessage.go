package gosocket

// Message describes an websocket message object.
type Message interface {

	// Data returns raw data bytes
	Data() []byte

	// Type returns type of Data
	Type() string

	// opCode returns opcode of message
	opCode() byte
}


type BinMsg interface {
	Message
}

type TextMsg interface {
	Message
}


// raw binary message
type rawMessage struct {
	opcode byte
	data []byte
}

func (msg *rawMessage) opCode() byte {
	return msg.opcode
}

func (msg *rawMessage) Data() []byte {
	return msg.data
}

func (msg *rawMessage) Type() string {
	switch (msg.opcode) {
		case M_TXT:
			return "text"
		case M_BIN:
			return "binary"
		case M_CLS:
			return "close"
		case M_PING:
			return "ping"
		case M_PONG:
			return "pong"
	}
	return "invalid"
}


// text message
type textMessage struct {
	*rawMessage
}


// Message describes CloseMsg websocket message.
type CloseMsg interface {
	Message

	// returns Closing code
	Code() int

	// returns Closing message
	Msg() string
}

// close message
type closeMessage struct {
	*rawMessage
}

func (msg *closeMessage) Code() int {
	var code int
	code = 0
	if len(msg.data) == 0 {
		return 1005
	}
	code |= (int(msg.data[0]) << 8)
	code |= (int(msg.data[1]))
	return code
}

func (msg *closeMessage) Msg() string {
	if len(msg.data) > 2 {
		return string(msg.data[2:])
	} else {
		return ""
	}
}

func NewTextMsg(data string) Message {
	return &textMessage{
		rawMessage: &rawMessage{
			data: []byte(data),
			opcode: M_TXT,
		},
	}
}

func NewBinMsg(data []byte) Message {
	return &rawMessage{
		data: data,
		opcode: M_BIN,
	}
}

func newCloseMsg(data []byte) CloseMsg {
	return &closeMessage{
		rawMessage: &rawMessage{
			data: data,
			opcode: M_CLS,
		},
	}
}

func NewCloseMsg(code int, str string) CloseMsg {
	var data []byte
	data = append(data, byte(code>>8), byte(code))
	data = append(data, []byte(str)...)
	return newCloseMsg(data)
}
