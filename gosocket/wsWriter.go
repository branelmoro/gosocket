package gosocket

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type wsWriter struct {
	*wsConn
	_msgLock sync.Mutex
	_frameLock sync.Mutex
}

func (w *wsWriter) isCloseInitiated() bool {
	return w.isCloseSent() || w.isCloseReceived() || w.isConnClosed()
}

func (w *wsWriter) getWriteFrame(op byte) *wFrame {
	return &wFrame{
		opcode: op,
		rsv1: false,
		rsv2: false,
		rsv3: false,
		isMasked: !w._isClient,
	}
}

func (w *wsWriter) ping() error {

	pingBytes := []byte{}
	frame := w.getWriteFrame(M_PING)
	frame.fin = true
	frame.data = pingBytes

	err := w.sendFrameBytes(frame.toBytes())

	if err == nil {
		w.pingSent = pingBytes
		w.setConnStatus(PING_SENT)
	} else {
		// close underlying tcp connection
		w.closeConn()
	}
	return err
}

func (w *wsWriter) pong() error {

	frame := w.getWriteFrame(M_PONG)
	frame.fin = true
	frame.data = w.pingReceived

	err := w.sendFrameBytes(frame.toBytes())

	if err == nil {
		w.setConnStatus(PONG_SENT)
		w.pingReceived = []byte{}
	} else {
		// close underlying tcp connection
		w.closeConn()
	}
	return err
}

func (w *wsWriter) Close(msg CloseMsg) error {

	w._frameLock.Lock()
	defer w._frameLock.Unlock()

	if w.isCloseSent() || w.isConnClosed() {
		return nil
	}

	if msg == nil {
		msg = NewCloseMsg(CC_NORMAL_CLOSURE, "normal close")
	}

	frame := w.getWriteFrame(M_CLS)
	frame.fin = true
	frame.data = msg.Data()

	// send close frame
	err := w.sendBytes(frame.toBytes())

	if err == nil {
		// mark close frame SendText``
		w.setConnStatus(CLOSE_SENT)
		if w.isCloseReceived() {
			// close underlying tcp connection
			w.closeConn()
		} else {
			// start reading incoming close frame
			reader := w.reader()
			reader.isClosing = true
			reader.start()
		}
	} else {
		// close underlying tcp connection
		w.closeConn()
	}
	return err
}

func (w *wsWriter) Send(msg Message) error {
	// return w.sendData(msg.opCode(), msg.Data())
	return nil
}

func (w *wsWriter) SendBinary(data []byte) error {
	return w.sendData(M_BIN, data)
}

func (w *wsWriter) SendText(str string) error {
	return w.sendData(M_TXT, []byte(str))
}

func (w *wsWriter) sendData(opcode byte, data []byte) error {
	var(
		err error
		frames [][]byte
	)
	
	frame := w.getWriteFrame(opcode)
	if w.isPerMessageDeflateEnabled() {
		data, err = w.flate.compress(data)
		if err != nil {
			fmt.Println("Err in data compression---", err)
		}
		data = data[:len(data)-4]
		frame.rsv1 = true
	}

	length := len(data)
	startIndex := 0
	for {
		frame.fin = length <= w.server.wsMaxFrameSize
		if frame.fin {
			frame.data = data[startIndex:]
			frames = append(frames, frame.toBytes())
			break
		} else {
			frame.data = data[startIndex : startIndex+w.server.wsMaxFrameSize]
			frames = append(frames, frame.toBytes())
		}
		frame.rsv1 = false
		startIndex += w.server.wsMaxFrameSize
		length -= w.server.wsMaxFrameSize
	}

	w._msgLock.Lock()
	defer w._msgLock.Unlock()
	for _, frameBytes := range frames {
		err = w.sendFrameBytes(frameBytes)
		if err != nil {
			// error in sending frame, close underlying tcp connection
			w.closeConn()
			return err
		}
	}
	return err
}

func (w *wsWriter) sendFrameBytes(frameData []byte) error {
	w._frameLock.Lock()
	defer w._frameLock.Unlock()
	if w.isCloseInitiated() {
		// connection has been closed or connection close has been intiated
		return newConnectionClosedError(w)
	}
	return w.sendBytes(frameData)
}

func (w *wsWriter) sendBytes(data []byte) error {

	// mark message writing start
	w.setConnStatus(WRITING_ON)
	defer w.setConnStatus(WRITING_OFF)
	defer w.setWriteTimeOut(time.Time{})

	var(
		numBytes int
		cntBytes int
		err error
	)

	sec := 10
	minBytes := sec * w.server.minIOSpeed()

	err = w.setWriteTimeOut(time.Now().Add(time.Duration(sec) * time.Second))

	if err != nil {
		return newSetWriteTimeoutError(err)
	}

	startIndex := 0

	cntBytes = 0

	for startIndex != len(data) {

		numBytes, err = w.write(data[startIndex:])

		cntBytes += numBytes

		if e, ok := err.(net.Error); ok && e.Timeout() {
			// timeout occured
			if cntBytes < minBytes {
				// return error, connection accept less data (numbytes bytes data) for 10 second
				// expecting minBytes (w.server.minIOSpeed() per second)
				return newSlowDataWriteError(cntBytes, sec)
			}
			err = w.setWriteTimeOut(time.Now().Add(time.Duration(sec) * time.Second))
			if err != nil {
				return newSetWriteTimeoutError(err)
			}
			err = nil
			cntBytes = 0
		}

		if err != nil{
			break
		}

		startIndex += numBytes
	}

	return err
}
