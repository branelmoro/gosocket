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

func newHttpMalformedError(str string) error {
	fmt.Println(str)
	return &httpError{
		error: fmt.Errorf("ERR_HTTP_MALFORMED: Http Malformed Error..." + str),
		_code: ERR_HTTP_MALFORMED,
	}
}

func newHttpUriError(err error) error {
	return &httpError{
		error: err,
		_code: ERR_HTTP_URI,
	}
}