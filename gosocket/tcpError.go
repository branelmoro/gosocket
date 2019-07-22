package gosocket

import(
	"io"
	"net"
)

type tcpError struct {
	error
	_code byte
}

func (e *tcpError) Code() byte {
	return e._code
}

func (e *tcpError) Detail() string {
	return "TCP Error:"
}

func isTimeout(err error) bool {
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}

func newTCPError(code byte, err error) error {
	return &tcpError{
		error: err,
		_code: code,
	}
}

func newMsgStartError(err error) error {
	if err == io.EOF {
		return newTCPError(ERR_TCP_CLOSE, err)
	} else {
		if isTimeout(err) {
			return newTCPError(ERR_INVALID, err)
		} else {
			return newTCPError(ERR_TCP_READ, err)
		}
	}
}

func newWriteError(err error) error {
	if err == io.EOF {
		return newTCPError(ERR_TCP_CLOSE, err)
	} else {
		return newTCPError(ERR_TCP_WRITE, err)
	}
}

func newReadError(err error) error {
	if err == io.EOF {
		return newTCPError(ERR_TCP_CLOSE, err)
	} else {
		return newTCPError(ERR_TCP_READ, err)
	}
}

func newSetReadTimeoutError(err error) error {
	if err == io.EOF {
		return newTCPError(ERR_TCP_CLOSE, err)
	} else {
		return newTCPError(ERR_SET_READTIMEOUT, err)
	}
}

func newSetWriteTimeoutError(err error) error {
	if err == io.EOF {
		return newTCPError(ERR_TCP_CLOSE, err)
	} else {
		return newTCPError(ERR_SET_WRITETIMEOUT, err)
	}
}
