package gosocket

import(
	"fmt"
)

const(
	ERR_HTTP_MALFORMED byte = 0
	ERR_HTTP_URI byte = 1
)

type httpError struct {
	error
	_code byte
}

func newHttpMalformedError() error {
	return &httpError{
		error: fmt.Errorf("ERR_HTTP_MALFORMED: Http Malformed Error."),
		_code: ERR_HTTP_MALFORMED,
	}
}

func newHttpUriError(err error) error {
	return &httpError{
		error: err,
		_code: ERR_HTTP_URI,
	}
}