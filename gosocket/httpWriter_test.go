package gosocket

// import (
// 	"fmt"
// 	"reflect"
// 	"testing"
// 	"time"
// )

// // -----------------httpWriter.Write tests start--------------------- //

// // set timeout failed at 1st place
// type testConnWriter1 struct {
// 	err error
// }
// func (c *testConnWriter1) setWriteTimeOut(t time.Time) error {
// 	return c.err
// }
// func (c *testConnWriter1) write(data []byte) (int, error) {
// 	return len(data), nil
// }
// func (c *testConnWriter1) minSpeed() int {
// 	return 1
// }
// func (c *testConnWriter1) close() error {
// 	return nil
// }
// func TestWriteTimeOutFailed1(t *testing.T) {
// 	data := []byte("this is data to write")
// 	err1 := fmt.Errorf("setWriteTimeOut failed error")
// 	writer := &testConnWriter1{
// 		err: err1,
// 	}
// 	w := &httpWriter {
// 		writer: writer,
// 		req: nil,
// 	}
// 	err := w.Write(data)
// 	if !reflect.DeepEqual(newSetWriteTimeoutError(err1), err) {
// 		t.Errorf("`HttpWriter.write` didn't timeout at first place... Error is - `%v`", err)
// 	}
// }


// // CustomNetError implements net.Error interface
// type testCustomNetError struct {
// 	error
// }
// func (c *testCustomNetError) Timeout() bool {
// 	return true
// }
// func (c *testCustomNetError) Temporary() bool {
// 	return true
// }


// // set timeout failed at 2nd place
// type testConnWriter2 struct {
// 	err error
// 	cnt int
// }
// func (c *testConnWriter2) setWriteTimeOut(t time.Time) error {
// 	if(c.cnt > 5) {
// 		return c.err
// 	}
// 	return nil
// }
// func (c *testConnWriter2) write(data []byte) (int, error) {
// 	if(c.cnt > 11) {
// 		return c.cnt, c.err
// 	}
// 	inc := len(data)
// 	if(inc >= 2) {
// 		inc = 2
// 	}
// 	c.cnt = c.cnt + inc
// 	return inc, nil
// }
// func (c *testConnWriter2) minSpeed() int {
// 	return 1
// }
// func (c *testConnWriter2) close() error {
// 	return nil
// }
// func TestWriteTimeOutFailed2(t *testing.T) {
// 	data := []byte("this is data to write")
// 	err1 := &testCustomNetError {
// 		error: fmt.Errorf("This is custom timeout net.Error"),
// 	}
// 	writer := &testConnWriter2{
// 		err: err1,
// 	}
// 	w := &httpWriter {
// 		writer: writer,
// 		req: nil,
// 	}
// 	err := w.Write(data)
// 	if !reflect.DeepEqual(newSetWriteTimeoutError(err1), err) {
// 		t.Errorf("`HttpWriter.write` didn't timeout at second place... Error is - `%v`", err)
// 	}
// }


// // check slow write error
// type testConnWriter3 struct {
// 	err error
// 	cnt int
// }
// func (c *testConnWriter3) setWriteTimeOut(t time.Time) error {
// 	return nil
// }
// func (c *testConnWriter3) write(data []byte) (int, error) {
// 	if(c.cnt > 5) {
// 		return c.cnt, c.err
// 	}
// 	inc := len(data)
// 	if(inc >= 2) {
// 		inc = 2
// 	}
// 	c.cnt = c.cnt + inc
// 	return inc, nil
// }
// func (c *testConnWriter3) minSpeed() int {
// 	return 1
// }
// func (c *testConnWriter3) close() error {
// 	return nil
// }
// func TestSlowWriteError(t *testing.T) {
// 	data := []byte("this is data to write")
// 	err1 := &testCustomNetError {
// 		error: fmt.Errorf("This is custom timeout net.Error"),
// 	}
// 	writer := &testConnWriter3{
// 		err: err1,
// 	}
// 	w := &httpWriter {
// 		writer: writer,
// 		req: nil,
// 	}
// 	err := w.Write(data)
// 	if !reflect.DeepEqual(newSlowDataWriteError(6, 10), err) {
// 		t.Errorf("`HttpWriter.write` didn't return slow write error... Error is - `%v`", err)
// 	}
// }


// // check write error
// type TestConnWriterError struct {
// 	err error
// }
// func (c *TestConnWriterError) setWriteTimeOut(t time.Time) error {
// 	return nil
// }
// func (c *TestConnWriterError) write(data []byte) (int, error) {
// 	return 0, c.err
// }
// func (c *TestConnWriterError) minSpeed() int {
// 	return 1
// }
// func (c *TestConnWriterError) close() error {
// 	return nil
// }
// func TestWriterError(t *testing.T) {
// 	data := []byte("this is data to write")
// 	err1 := fmt.Errorf("This is custom write error")
// 	writer := &TestConnWriterError{
// 		err: err1,
// 	}
// 	w := &httpWriter {
// 		writer: writer,
// 		req: nil,
// 	}
// 	err := w.Write(data)
// 	if !reflect.DeepEqual(newWriteError(err1), err) {
// 		t.Errorf("`HttpWriter.write` didn't return write error... Error is - `%v`", err)
// 	}
// }


// // normal behaviour
// type testConnWriter struct {}
// func (c *testConnWriter) setWriteTimeOut(t time.Time) error {
// 	return nil
// }
// func (c *testConnWriter) write(data []byte) (int, error) {
// 	return len(data), nil
// }
// func (c *testConnWriter) minSpeed() int {
// 	return 1
// }
// func (c *testConnWriter) close() error {
// 	return nil
// }
// func TestWrite(t *testing.T) {
// 	data := []byte("this is data to write")
// 	writer := &testConnWriter{}
// 	w := &httpWriter {
// 		writer: writer,
// 		req: nil,
// 	}
// 	err := w.Write(data)
// 	if err != nil {
// 		t.Errorf("`HttpWriter.write` is not working as expected... Error is - `%v`", err)
// 	}
// }

// // -----------------httpWriter.Write tests end--------------------- //


// // -----------------httpWriter.UpgradeToWebsocket tests start--------------------- //

// // websocket request
// type testHttpReq struct {}
// func (c *testHttpReq) Header() map[string]string {
// 	return map[string]string {
// 		"upgrade": "websocket",
// 		"connection": "Upgrade",
// 		"Sec-WebSocket-Accept": "s3pPLMBiTxaQ9kYGzzhZRbK+xOo=",
// 		"Sec-WebSocket-Version": "13",
// 		"sec-websocket-extensions": "permessage-deflate",
// 		"header1": "abc",
// 		"header2": "abcd",
// 	}
// }
// func (c *testHttpReq) isWebSocketRequest() bool {
// 	return true
// }
// func TestUpgradeToWebsocket(t *testing.T) {
// 	writer := &testConnWriter{}
// 	w := &httpWriter {
// 		writer: writer,
// 		req: &testHttpReq{},
// 	}
// 	err := w.UpgradeToWebsocket(nil)
// 	if err != nil {
// 		t.Errorf("`HttpWriter.UpgradeToWebsocket` does not upgrate websocket request")
// 	}
// }


// // websocket request with write error
// func TestUpgradeToWebsocketWriteError(t *testing.T) {
// 	writer := &testConnWriter1{
// 		err: fmt.Errorf("setWriteTimeOut failed error"),
// 	}
// 	w := &httpWriter {
// 		writer: writer,
// 		req: &testHttpReq{},
// 	}
// 	err := w.UpgradeToWebsocket(nil)
// 	if err == nil {
// 		t.Errorf("`HttpWriter.UpgradeToWebsocket` does return write error")
// 	}
// }


// // non websocket request
// type testHttpReq1 struct {}
// func (c *testHttpReq1) Header() map[string]string {
// 	return nil
// }
// func (c *testHttpReq1) isWebSocketRequest() bool {
// 	return false
// }
// func TestUpgradeToWebsocket1(t *testing.T) {
// 	writer := &testConnWriter{}
// 	w := &httpWriter {
// 		writer: writer,
// 		req: &testHttpReq1{},
// 	}
// 	err := w.UpgradeToWebsocket(nil)
// 	if err == nil {
// 		t.Errorf("`HttpWriter.UpgradeToWebsocket` does not return error for non websocket request")
// 	}
// }

// // -----------------httpWriter.UpgradeToWebsocket tests end--------------------- //