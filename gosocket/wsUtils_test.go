package gosocket

import (
	"reflect"
	"testing"
)

func TestGenerateSecWebSocketAccept(t *testing.T) {
	inputWsAccept := "dGhlIHNhbXBsZSBub25jZQ=="
	expectedWsAccept := "s3pPLMBiTxaQ9kYGzzhZRbK+xOo="
	outputWsAccept := generateSecWebSocketAccept(inputWsAccept)
    if outputWsAccept != expectedWsAccept {
       t.Errorf("`generateSecWebSocketAccept` does not return expected Web-Socket-Accept - `%s` for input Web-Socket-Accept `%s`", expectedWsAccept, expectedWsAccept)
    }
}

func TestGenerateWsUpgradeHeader(t *testing.T) {
	inputWsAccept := "dGhlIHNhbXBsZSBub25jZQ=="
	expectedWsAccept := "s3pPLMBiTxaQ9kYGzzhZRbK+xOo="
	requestHeader := make(map[string]string)
	requestHeader["sec-websocket-key"] = inputWsAccept

	responseHeader := generateWsUpgradeHeader(requestHeader, nil, nil)

	expectedHeader := map[string]string {
		"upgrade": "websocket",
		"connection": "Upgrade",
		"Sec-WebSocket-Accept": expectedWsAccept,
		"Sec-WebSocket-Version": "13",
	}

	if !reflect.DeepEqual(responseHeader, expectedHeader) {
		t.Errorf("`generateWsUpgradeHeader` not working as expected... Expecting responseHeader - %v for requestHeader - %v but returned responseHeader is %v", expectedHeader, requestHeader, responseHeader)
	}
}

func TestGenerateWsUpgradeHeaderWithDeflate(t *testing.T) {
	inputWsAccept := "dGhlIHNhbXBsZSBub25jZQ=="
	expectedWsAccept := "s3pPLMBiTxaQ9kYGzzhZRbK+xOo="
	requestHeader := make(map[string]string)
	requestHeader["sec-websocket-key"] = inputWsAccept
	requestHeader["sec-websocket-extensions"] = "permessage-deflate"

	responseHeader := generateWsUpgradeHeader(requestHeader, nil, nil)

	expectedHeader := map[string]string {
		"upgrade": "websocket",
		"connection": "Upgrade",
		"Sec-WebSocket-Accept": expectedWsAccept,
		"Sec-WebSocket-Version": "13",
		"sec-websocket-extensions": "permessage-deflate",
	}

	if !reflect.DeepEqual(responseHeader, expectedHeader) {
		t.Errorf("`generateWsUpgradeHeader` not working as expected... Expecting responseHeader - %v for requestHeader - %v but returned responseHeader is %v", expectedHeader, requestHeader, responseHeader)
	}
}

func TestGenerateWsUpgradeHeaderWithOptions(t *testing.T) {
	inputWsAccept := "dGhlIHNhbXBsZSBub25jZQ=="
	expectedWsAccept := "s3pPLMBiTxaQ9kYGzzhZRbK+xOo="
	requestHeader := make(map[string]string)
	requestHeader["sec-websocket-key"] = inputWsAccept

	options := &WsOptions{
		Headers: map[string]string {
			"header1": "abc",
			"header2": "abcd",
		},
		WsData: nil,
	}

	responseHeader := generateWsUpgradeHeader(requestHeader, options, nil)

	expectedHeader := map[string]string {
		"upgrade": "websocket",
		"connection": "Upgrade",
		"Sec-WebSocket-Accept": expectedWsAccept,
		"Sec-WebSocket-Version": "13",
		"header1": "abc",
		"header2": "abcd",
	}

	if !reflect.DeepEqual(responseHeader, expectedHeader) {
		t.Errorf("`generateWsUpgradeHeader` not working as expected... Expecting responseHeader - %v for requestHeader - %v but returned responseHeader is %v", expectedHeader, requestHeader, responseHeader)
	}
}

func TestGenerateWsUpgradeHeaderWithDeflateAndCustomHeader(t *testing.T) {
	inputWsAccept := "dGhlIHNhbXBsZSBub25jZQ=="
	expectedWsAccept := "s3pPLMBiTxaQ9kYGzzhZRbK+xOo="
	requestHeader := make(map[string]string)
	requestHeader["sec-websocket-key"] = inputWsAccept
	requestHeader["sec-websocket-extensions"] = "permessage-deflate"

	options := &WsOptions{
		Headers: map[string]string {
			"header1": "abc",
			"header2": "abcd",
		},
		WsData: nil,
	}

	responseHeader := generateWsUpgradeHeader(requestHeader, options, nil)

	expectedHeader := map[string]string {
		"upgrade": "websocket",
		"connection": "Upgrade",
		"Sec-WebSocket-Accept": expectedWsAccept,
		"Sec-WebSocket-Version": "13",
		"sec-websocket-extensions": "permessage-deflate",
		"header1": "abc",
		"header2": "abcd",
	}

	if !reflect.DeepEqual(responseHeader, expectedHeader) {
		t.Errorf("`generateWsUpgradeHeader` not working as expected... Expecting responseHeader - %v for requestHeader - %v but returned responseHeader is %v", expectedHeader, requestHeader, responseHeader)
	}
}
