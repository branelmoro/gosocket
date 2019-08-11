package gosocket

import (
	"encoding/hex"
	"fmt"
)

type Error interface {
	error
	Code() byte
	Detail() string
}

type wsFrameError struct {
	error
	_code byte
	frame *wsFrame
}

func (e *wsFrameError) Code() byte {
	return e._code
}

func (e *wsFrameError) Detail() string {
	return hex.Dump(e.frame.toBytes())
}

func newUnidentifiedFrameError(frame *wsFrame) error {
	return &wsFrameError{
		error: fmt.Errorf("ERR_UNIDENTIFIED_FRAME: Unidentified frame received with Opcode - 0x%X.", frame.opcode()),
		frame: frame,
		_code: ERR_UNIDENTIFIED_FRAME,
	}
}

func newInvalidMessageStartError(frame *wsFrame) error {
	return &wsFrameError{
		error: fmt.Errorf("ERR_INVALID_MESSAGE_START: Continuation frame received with opcode - 0x%X at start of message.", frame.opcode()),
		frame: frame,
		_code: ERR_INVALID_MESSAGE_START,
	}
}

func newExpectingCloseFrameError(frame *wsFrame) error {
	return &wsFrameError{
		error: fmt.Errorf("ERR_EXPECTING_CLOSE_FRAME: Expecting close frame.. but %s Frame received with opcode - 0x%X after sending close frame.", frame.getType(), frame.opcode()),
		frame: frame,
		_code: ERR_EXPECTING_CLOSE_FRAME,
	}
}

func newExpectingContinueFrameError(frame *wsFrame) error {
	return &wsFrameError{
		error: fmt.Errorf("ERR_EXPECTING_CONTINUE_FRAME: Expecting continue frame.. but %s frame received with opcode - 0x%X in continued message", frame.getType(), frame.opcode()),
		frame: frame,
		_code: ERR_EXPECTING_CONTINUE_FRAME,
	}
}

func newExpectingMaskedFrameError(frame *wsFrame) error {
	return &wsFrameError{
		error: fmt.Errorf("ERR_EXPECTING_MASKED_FRAME: Expecting masked frame.. Received unmasked frame from client."),
		frame: frame,
		_code: ERR_EXPECTING_MASKED_FRAME,
	}
}

func newExpectingUnmaskedFrameError(frame *wsFrame) error {
	return &wsFrameError{
		error: fmt.Errorf("ERR_EXPECTING_UNMASKED_FRAME: Expecting unmasked frame.. Received masked frame from server."),
		frame: frame,
		_code: ERR_EXPECTING_UNMASKED_FRAME,
	}
}

func newFramePayloadLengthError(frame *wsFrame, isSmall bool) error {
	e := wsFrameError{
		frame: frame,
		_code: ERR_FRAME_PAYLOAD_LENGTH,
	}
	if isSmall {
		e.error = fmt.Errorf("ERR_FRAME_PAYLOAD_LENGTH: Less payload length - %d bytes found in 16 bit length bytes.. expecting more than or equal to 126.", frame.payloadLength())
	} else {
		e.error = fmt.Errorf("ERR_FRAME_PAYLOAD_LENGTH: Less payload length - %d bytes found in 64 bit length bytes.. expecting more than or equal to 65537.", frame.payloadLength())
	}
	return &e
}

func newCloseFrameLengthError(frame *wsFrame) error {
	return &wsFrameError{
		error: fmt.Errorf("ERR_CLOSE_FRAME_LENGTH: Close frame's payload length can't be one byte.. Close Code needs two bytes."),
		frame: frame,
		_code: ERR_CLOSE_FRAME_LENGTH,
	}
}

func newEmptyDataFrameError(frame *wsFrame) error {
	return &wsFrameError{
		error: fmt.Errorf("ERR_EMPTY_DATA_FRAME: Data frame's payload length can't be zero byte.. expecting more than or equal to 1 bytes."),
		frame: frame,
		_code: ERR_EMPTY_DATA_FRAME,
	}
}

func newControlFrameFinError(frame *wsFrame) error {
	return &wsFrameError{
		error: fmt.Errorf("ERR_CONTROL_FRAME_FIN: Control frame must have fin bit set to 1."),
		frame: frame,
		_code: ERR_CONTROL_FRAME_FIN,
	}
}

func newControlFrameRsv1Error(frame *wsFrame) error {
	return &wsFrameError{
		error: fmt.Errorf("ERR_CONTROL_FRAME_RSV1: Control frame must not have rsv1 bit set to 1."),
		frame: frame,
		_code: ERR_CONTROL_FRAME_RSV1,
	}
}

func newControlFrameLengthError(frame *wsFrame) error {
	return &wsFrameError{
		error: fmt.Errorf("ERR_CONTROL_FRAME_LENGTH: Control frame payload length must be less than or equal to 125."),
		frame: frame,
		_code: ERR_CONTROL_FRAME_LENGTH,
	}
}






type wsError struct {
	error
	msg []byte
	_code byte
}

func (e *wsError) Code() byte {
	return e._code
}

func (e *wsError) Detail() string {
	return ""
}

func newSlowDataReadError(cnt int, time int) error {
	return &wsError{
		error: fmt.Errorf("ERR_SLOW_DATA_READ: Endpoint is sending data at slow rate - %d bytes per %d seconds(approximately).", cnt, time),
		_code: ERR_SLOW_DATA_READ,
	}
}

func newSlowDataWriteError(cnt int, time int) error {
	return &wsError{
		error: fmt.Errorf("ERR_SLOW_DATA_WRITE: Endpoint is accepting data at slow rate - %d bytes per %d seconds(approximately).", cnt, time),
		_code: ERR_SLOW_DATA_WRITE,
	}
}

func newUtf8TextError(message []byte) error {
	return &wsError{
		error: fmt.Errorf("ERR_TEXT_UTF8: Invalid utf8 text received in text message."),
		_code: ERR_TEXT_UTF8,
		msg: message,
	}
}

func newTooBigMessageError(lenght int) error {
	return &wsError{
		error: fmt.Errorf("ERR_BIG_MESSAGE: Too big message received."),
		_code: ERR_BIG_MESSAGE,
	}
}

func newConnectionClosedError(w *wsWriter) error {
	return &wsError{
		error: fmt.Errorf("ERR_CONNECTION_CLOSED: Connection close had been initiated or connection is closed."),
		_code: ERR_CONNECTION_CLOSED,
	}
}

