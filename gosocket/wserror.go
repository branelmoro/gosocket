package gosocket

const (
	NORMAL_CLOSURE              = 1000
	GOING_AWAY                  = 1001
	PROTOCOL_ERROR              = 1002
	DATA_TYPE_NOT_ACCEPTABLE    = 1003
	NO_STATUS_CODE_PRESENT      = 1005
	CONNECTION_CLOSED_ABNORMALY = 1006
	NONCONSISTANT_MESSAGE_DATA  = 1007
	MESSAGE_POLICY_VIOLATION    = 1008
	TOO_BIG_MESSAGE             = 1009
	UNEXPECTED_ERROR            = 1011
	TLS_HANDSHAKE_FAILURE       = 1015
)

const (
	READ_ERROR              = 0x01
	OPCODE_ERROR            = 0x02
	MASK_BIT_ERROR          = 0x03
	PAYLOAD_LENGTH_ERROR    = 0x04
	MESSAGE_LENGTH_ERROR    = 0x05
	EOF_ERROR               = 0x06
)

type WsError struct {
	message string
	code byte
}

func (e *WsError) Error() string {
    return e.message
}

func NewWsError(code byte, message string) error {
    return &WsError{
        message: message,
        code: code,
    }
}
