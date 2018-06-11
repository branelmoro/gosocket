package gosocket

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
