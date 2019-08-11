package gosocket

const (
	HttpContinue           int =  100 // RFC 7231, 6.2.1
	HttpSwitchingProtocols int =  101 // RFC 7231, 6.2.2
	HttpProcessing         int =  102 // RFC 2518, 10.1
	HttpEarlyHints         int =  103 // RFC 8297

	HttpOK                   int =  200 // RFC 7231, 6.3.1
	HttpCreated              int =  201 // RFC 7231, 6.3.2
	HttpAccepted             int =  202 // RFC 7231, 6.3.3
	HttpNonAuthoritativeInfo int =  203 // RFC 7231, 6.3.4
	HttpNoContent            int =  204 // RFC 7231, 6.3.5
	HttpResetContent         int =  205 // RFC 7231, 6.3.6
	HttpPartialContent       int =  206 // RFC 7233, 4.1
	HttpMultiStatus          int =  207 // RFC 4918, 11.1
	HttpAlreadyReported      int =  208 // RFC 5842, 7.1
	HttpIMUsed               int =  226 // RFC 3229, 10.4.1

	HttpMultipleChoices   int =  300 // RFC 7231, 6.4.1
	HttpMovedPermanently  int =  301 // RFC 7231, 6.4.2
	HttpFound             int =  302 // RFC 7231, 6.4.3
	HttpSeeOther          int =  303 // RFC 7231, 6.4.4
	HttpNotModified       int =  304 // RFC 7232, 4.1
	HttpUseProxy          int =  305 // RFC 7231, 6.4.5
	_                       int =  306 // RFC 7231, 6.4.6 (Unused)
	HttpTemporaryRedirect int =  307 // RFC 7231, 6.4.7
	HttpPermanentRedirect int =  308 // RFC 7538, 3

	HttpBadRequest                   int =  400 // RFC 7231, 6.5.1
	HttpUnauthorized                 int =  401 // RFC 7235, 3.1
	HttpPaymentRequired              int =  402 // RFC 7231, 6.5.2
	HttpForbidden                    int =  403 // RFC 7231, 6.5.3
	HttpNotFound                     int =  404 // RFC 7231, 6.5.4
	HttpMethodNotAllowed             int =  405 // RFC 7231, 6.5.5
	HttpNotAcceptable                int =  406 // RFC 7231, 6.5.6
	HttpProxyAuthRequired            int =  407 // RFC 7235, 3.2
	HttpRequestTimeout               int =  408 // RFC 7231, 6.5.7
	HttpConflict                     int =  409 // RFC 7231, 6.5.8
	HttpGone                         int =  410 // RFC 7231, 6.5.9
	HttpLengthRequired               int =  411 // RFC 7231, 6.5.10
	HttpPreconditionFailed           int =  412 // RFC 7232, 4.2
	HttpRequestEntityTooLarge        int =  413 // RFC 7231, 6.5.11
	HttpRequestURITooLong            int =  414 // RFC 7231, 6.5.12
	HttpUnsupportedMediaType         int =  415 // RFC 7231, 6.5.13
	HttpRequestedRangeNotSatisfiable int =  416 // RFC 7233, 4.4
	HttpExpectationFailed            int =  417 // RFC 7231, 6.5.14
	HttpTeapot                       int =  418 // RFC 7168, 2.3.3
	HttpMisdirectedRequest           int =  421 // RFC 7540, 9.1.2
	HttpUnprocessableEntity          int =  422 // RFC 4918, 11.2
	HttpLocked                       int =  423 // RFC 4918, 11.3
	HttpFailedDependency             int =  424 // RFC 4918, 11.4
	HttpTooEarly                     int =  425 // RFC 8470, 5.2.
	HttpUpgradeRequired              int =  426 // RFC 7231, 6.5.15
	HttpPreconditionRequired         int =  428 // RFC 6585, 3
	HttpTooManyRequests              int =  429 // RFC 6585, 4
	HttpRequestHeaderFieldsTooLarge  int =  431 // RFC 6585, 5
	HttpUnavailableForLegalReasons   int =  451 // RFC 7725, 3

	HttpInternalServerError           int =  500 // RFC 7231, 6.6.1
	HttpNotImplemented                int =  501 // RFC 7231, 6.6.2
	HttpBadGateway                    int =  502 // RFC 7231, 6.6.3
	HttpServiceUnavailable            int =  503 // RFC 7231, 6.6.4
	HttpGatewayTimeout                int =  504 // RFC 7231, 6.6.5
	HttpHTTPVersionNotSupported       int =  505 // RFC 7231, 6.6.6
	HttpVariantAlsoNegotiates         int =  506 // RFC 2295, 8.1
	HttpInsufficientStorage           int =  507 // RFC 4918, 11.5
	HttpLoopDetected                  int =  508 // RFC 5842, 7.2
	HttpNotExtended                   int =  510 // RFC 2774, 7
	HttpNetworkAuthenticationRequired int =  511 // RFC 6585, 6
)

var httpStatusText =  map[int]string{
	HttpContinue:           "Continue",
	HttpSwitchingProtocols: "Switching Protocols",
	HttpProcessing:         "Processing",
	HttpEarlyHints:         "Early Hints",

	HttpOK:                   "OK",
	HttpCreated:              "Created",
	HttpAccepted:             "Accepted",
	HttpNonAuthoritativeInfo: "Non-Authoritative Information",
	HttpNoContent:            "No Content",
	HttpResetContent:         "Reset Content",
	HttpPartialContent:       "Partial Content",
	HttpMultiStatus:          "Multi-Status",
	HttpAlreadyReported:      "Already Reported",
	HttpIMUsed:               "IM Used",

	HttpMultipleChoices:   "Multiple Choices",
	HttpMovedPermanently:  "Moved Permanently",
	HttpFound:             "Found",
	HttpSeeOther:          "See Other",
	HttpNotModified:       "Not Modified",
	HttpUseProxy:          "Use Proxy",
	HttpTemporaryRedirect: "Temporary Redirect",
	HttpPermanentRedirect: "Permanent Redirect",

	HttpBadRequest:                   "Bad Request",
	HttpUnauthorized:                 "Unauthorized",
	HttpPaymentRequired:              "Payment Required",
	HttpForbidden:                    "Forbidden",
	HttpNotFound:                     "Not Found",
	HttpMethodNotAllowed:             "Method Not Allowed",
	HttpNotAcceptable:                "Not Acceptable",
	HttpProxyAuthRequired:            "Proxy Authentication Required",
	HttpRequestTimeout:               "Request Timeout",
	HttpConflict:                     "Conflict",
	HttpGone:                         "Gone",
	HttpLengthRequired:               "Length Required",
	HttpPreconditionFailed:           "Precondition Failed",
	HttpRequestEntityTooLarge:        "Request Entity Too Large",
	HttpRequestURITooLong:            "Request URI Too Long",
	HttpUnsupportedMediaType:         "Unsupported Media Type",
	HttpRequestedRangeNotSatisfiable: "Requested Range Not Satisfiable",
	HttpExpectationFailed:            "Expectation Failed",
	HttpTeapot:                       "I'm a teapot",
	HttpMisdirectedRequest:           "Misdirected Request",
	HttpUnprocessableEntity:          "Unprocessable Entity",
	HttpLocked:                       "Locked",
	HttpFailedDependency:             "Failed Dependency",
	HttpTooEarly:                     "Too Early",
	HttpUpgradeRequired:              "Upgrade Required",
	HttpPreconditionRequired:         "Precondition Required",
	HttpTooManyRequests:              "Too Many Requests",
	HttpRequestHeaderFieldsTooLarge:  "Request Header Fields Too Large",
	HttpUnavailableForLegalReasons:   "Unavailable For Legal Reasons",

	HttpInternalServerError:           "Internal Server Error",
	HttpNotImplemented:                "Not Implemented",
	HttpBadGateway:                    "Bad Gateway",
	HttpServiceUnavailable:            "Service Unavailable",
	HttpGatewayTimeout:                "Gateway Timeout",
	HttpHTTPVersionNotSupported:       "HTTP Version Not Supported",
	HttpVariantAlsoNegotiates:         "Variant Also Negotiates",
	HttpInsufficientStorage:           "Insufficient Storage",
	HttpLoopDetected:                  "Loop Detected",
	HttpNotExtended:                   "Not Extended",
	HttpNetworkAuthenticationRequired: "Network Authentication Required",
}


// close codes
const (
	CC_NORMAL_CLOSURE int = iota + 1000
	CC_GOING_AWAY
	CC_PROTOCOL_ERROR
	CC_UNACCEPTABLE_DATA
	_
	CC_NO_STATUS_CODE // MUST NOT be set as a status code
	CC_ABNORMAL_CLOSE // MUST NOT be set as a status code
	CC_INCONSISTANT_DATA
	CC_POLICY_VIOLATION
	CC_BIG_MESSAGE
	CC_EXTENSSION_MISSING
	CC_UNEXPECTED_ERROR
	CC_SERVICE_RESTART
	CC_TRY_AGAIN // temporary error e.g. server overloaded
	CC_BAD_GATEWAY // upstream server error
	CC_TLS_FAILURE // MUST NOT be set as a status code
)


// message codes
const (
	M_CONTINUE byte = 0x00
	M_TXT      byte = 0x01
	M_BIN      byte = 0x02
	M_CLS      byte = 0x08
	M_PING     byte = 0x09
	M_PONG     byte = 0x0a
)


const (
	// message read timeout error
	ERR_INVALID byte = iota

	// tcp errors
	ERR_TCP_READ
	ERR_TCP_WRITE
	ERR_TCP_CLOSE
	ERR_SET_READTIMEOUT
	ERR_SET_WRITETIMEOUT

	// ws protocol errors
	// ws frame erros
	ERR_UNIDENTIFIED_FRAME
	ERR_INVALID_MESSAGE_START
	ERR_EXPECTING_CLOSE_FRAME
	ERR_EXPECTING_CONTINUE_FRAME
	ERR_CONTROL_FRAME_FIN
	ERR_CONTROL_FRAME_RSV1
	ERR_EXPECTING_MASKED_FRAME
	ERR_EXPECTING_UNMASKED_FRAME
	ERR_EMPTY_DATA_FRAME
	ERR_CONTROL_FRAME_LENGTH
	ERR_FRAME_PAYLOAD_LENGTH
	ERR_CLOSE_FRAME_LENGTH

	ERR_TEXT_UTF8

	ERR_SLOW_DATA_READ
	ERR_SLOW_DATA_WRITE

	// ws message errors
	ERR_BIG_MESSAGE

	ERR_CONNECTION_CLOSED
)

