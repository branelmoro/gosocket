package gosocket

import (
	"fmt"
	"net"
	"reflect"
	"time"
	"unicode/utf8"
)

func isAllowedMessageLength(length int) bool {
	return length <= serverConf.wsMaxMessageSize
}

type wsReader struct {
	*wsConn
	isClosing bool
}

func (r *wsReader) processData(opcode byte, data []byte) {
	writer := r.writer()
	switch (opcode) {
		case M_TXT:
			OnText(writer, string(data))
			break
		case M_BIN:
			OnBinary(writer, data)
			break
	}
}

func (r *wsReader) processControl(opcode byte, data []byte) {
	writer := r.writer()
	switch (opcode) {
		case M_CLS:
			msg := newCloseMsg(data)
			r.setConnStatus(CLOSE_RECEIVED) // mark close frame received
			if r.isCloseSent() {
				r.closeConn() // close underlying tcp connection
			} else {
				// send close frame
				// if msg.Code() == CC_NORMAL_CLOSURE ||
				// 	msg.Code() == CC_GOING_AWAY ||
				// 	msg.Code() == CC_NO_STATUS_CODE {
				if true {
					// connection not closed yet, resend close frame
					err := writer.Close(msg)
					if err != nil {
						go OnError(writer, err)
					}
				}
			}
			go OnClose(writer, msg)
			break
		case M_PING:
			// ping
			r.pingReceived = data
			r.setConnStatus(PING_RECEIVED)
			go OnPing(writer)
			err := writer.pong()
			if err != nil {
				go OnError(writer, err)
			}
			break
		case M_PONG:
			// pong
			if r.connStatus&PING_SENT == PING_SENT {
				go OnPong(writer)
				r.setConnStatus(PONG_RECEIVED)
				if reflect.DeepEqual(r.pingSent, data) {
					r.pingSent = []byte{}
				} else {
					// invalid data received in pong
					// OnError(writer)
				}
			}
			break
	}
}

func (r *wsReader) handleError(err error) {
	writer := r.writer()
	e, _ := err.(Error)
	switch(e.Code()) {
		// message read timeout error
		case ERR_INVALID:
			// ignore this error
			if r.isClosing {
				// close frame didn't receive
				r.closeConn()
				go OnError(writer, err)
			} else {
				r.setTimeOut(time.Time{})
			}
			break

		// tcp errors
		case ERR_TCP_READ,
			ERR_TCP_WRITE,
			ERR_SET_READTIMEOUT,
			ERR_SET_WRITETIMEOUT:
			r.closeConn()
			go OnError(writer, err)
			go OnClose(writer, NewCloseMsg(CC_UNEXPECTED_ERROR, ""))
			break

		case ERR_TCP_CLOSE:
			fmt.Println("here-------------------------------")
			r.handleTCPClose(err)
			break

		// ws protocol errors
		// ws frame errors
		case ERR_UNIDENTIFIED_FRAME,
			ERR_INVALID_MESSAGE_START,
			ERR_EXPECTING_CONTINUE_FRAME,
			ERR_CONTROL_FRAME_FIN,
			ERR_CONTROL_FRAME_LENGTH,
			ERR_EXPECTING_MASKED_FRAME,
			ERR_EXPECTING_UNMASKED_FRAME:
			// these are protocol errors
			msg := NewCloseMsg(CC_PROTOCOL_ERROR, "")
			writer.Close(msg)
			go OnError(writer, err)
			go OnClose(writer, msg)
			break

		case ERR_EMPTY_DATA_FRAME,
			ERR_FRAME_PAYLOAD_LENGTH,
			ERR_CLOSE_FRAME_LENGTH:
			// these are protocol errors
			msg := NewCloseMsg(CC_POLICY_VIOLATION, "")
			writer.Close(msg)
			go OnError(writer, err)
			go OnClose(writer, msg)
			break

		case ERR_EXPECTING_CLOSE_FRAME:
			r.closeConn()
			go OnError(writer, err)
			break

		case ERR_SLOW_DATA_READ:
			msg := NewCloseMsg(CC_POLICY_VIOLATION, "")
			writer.Close(msg)
			r.closeConn()
			go OnError(writer, err)
			break

		case ERR_SLOW_DATA_WRITE:
			r.closeConn()
			go OnError(writer, err)
			break

		case ERR_TEXT_UTF8:
			msg := NewCloseMsg(CC_INCONSISTANT_DATA, "")
			writer.Close(msg)
			r.closeConn()
			go OnError(writer, err)
			break

		// ws message errors
		case ERR_BIG_MESSAGE:
			msg := NewCloseMsg(CC_BIG_MESSAGE, "")
			writer.Close(msg)
			r.closeConn()
			go OnError(writer, err)
			break
	}
}

func (r *wsReader) start() {
	r._readLock.Lock()
	defer r._readLock.Unlock()

	if r.isCloseReceived() || r.isConnClosed() {
		return
	}

	// mark message read started
	r.setConnStatus(READING_ON)

	var (
		frame *wsFrame
		opcode byte
		readBytes []byte
		isMsgFinished bool
		isCloseTimeoutSet bool
		messageLength int
		rsv1 bool
		err error
	)

	isMsgFinished = true
	messageLength = 0

	isCloseTimeoutSet = false

	for {
		if r.isClosing && !isCloseTimeoutSet {
			err = r.setReadTimeOut(time.Now().Add(serverConf.wsCloseReadTimeout * time.Second))
			if err != nil {
				// return error, Close frame not received in serverConf.wsCloseReadTimeout time
				r.handleError(newSetReadTimeoutError(err))
				break
			}
			isCloseTimeoutSet = true
		}
		frame, err = r.readFrameHeader(isMsgFinished)
		if err != nil {
			r.handleError(err)
			break
		}
		fmt.Println("Frame Opcode : ", frame.opcode(), ", Frame data-----", frame.toBytes())
		if frame.isControlFrame() {
			err = r.readFrameData(frame)
			if err != nil {
				r.handleError(err)
				break
			}
			controlPayload := frame.payload()
			if frame.isMasked() { frame.unMask(controlPayload) }
			r.processControl(frame.opcode(), controlPayload)
		} else {
			messageLength += frame.length()
			if !isAllowedMessageLength(messageLength) {
				// close r connection, message length more than acceptable
				r.handleError(newTooBigMessageError(messageLength))
				break
			}
			err = r.readFrameData(frame)
			if err != nil {
				r.handleError(err)
				break
			}
			if isMsgFinished {
				opcode = frame.opcode()
				rsv1 = frame.rsv1()
			}
			readBytes = append(readBytes, frame.payload()...)
			if frame.isMasked() { frame.unMask(readBytes) }
			isMsgFinished = frame.fin()
			if isMsgFinished {
				// message finished
				fmt.Println("Received Message data-----", readBytes)
				if r.isPerMessageDeflateEnabled() && rsv1 {
					// decompress message
					readBytes, err = r.flate.decompress(readBytes)
					if err != nil {
						fmt.Println("Err in decompress-----", err)
					}
					// readBytes = decompress(readBytes)
					fmt.Println("Decompressed Received Message data-----", readBytes)
				}

				if opcode == M_TXT && !utf8.Valid(readBytes) {
					// text message is not valid utf-8 string
					r.handleError(newUtf8TextError(readBytes))
					break
				}
				go r.processData(opcode, readBytes)
				messageLength = 0
				readBytes = []byte{}
			}
		}
		if r.isCloseReceived() {
			// close frame already received, dont read any more messages
			break
		}
		r.isClosing = r.isCloseSent()
	}

	// mark message reading done
	r.setConnStatus(READING_OFF)

	// if true {
	// 	// server is going Down
	// 	err = r.writer().Close(NewCloseMsg(GOING_AWAY, "Shutting Down"))
	// 	if err != nil {
	// 		// handle close error
	// 	}
	// }

}

func (r *wsReader) setTimeOut(t time.Time) error {
	if !r.isClosing {
		return r.setReadTimeOut(t)
	}
	return nil
}

func (r *wsReader) readFrameHeader(isFirst bool) (*wsFrame, error) {
	var(
		frame wsFrame
		numBytes int
		readBytes []byte
		err error
	)

	if isFirst {
		// set timeout on message read start
		err = r.setTimeOut(time.Now().Add(1 * time.Millisecond))
	} else {
		err = r.setTimeOut(time.Now().Add(serverConf.wsHeaderReadTimeout * time.Second))
	}
	if err != nil {
		// error in setting timeout
		return &frame, newSetReadTimeoutError(err)
	}

	// read first byte
	numBytes, readBytes, err = r.read(1)
	if err == nil {
		frame.firstByte = readBytes[0]
	} else {
		if isFirst {
			return &frame, newMsgStartError(err)
		} else {
			return &frame, newReadError(err)
		}
	}

	switch (frame.opcode()) {
		case M_CONTINUE, M_TXT, M_BIN, M_CLS, M_PING, M_PONG:
			break
		default:
			return &frame, newUnidentifiedFrameError(&frame)
	}

	if frame.isControlFrame() && !frame.fin() {
		// control frame must be final frame with fin bit set to 1
		return &frame, newControlFrameFinError(&frame)
	}

	if isFirst {
		if frame.opcode() == M_CONTINUE {
			// message continuation frame received at start of message
			return &frame, newInvalidMessageStartError(&frame)
		}
		if r.isClosing && frame.opcode() != M_CLS {
			// return error, didn't receive close frame after sending close frame
			return &frame, newExpectingCloseFrameError(&frame)
		}
	} else {
		if frame.opcode() != M_CONTINUE && !frame.isControlFrame() {
			// invalid opcode received in message continuation frame
			return &frame, newExpectingContinueFrameError(&frame)
		}
	}

	if isFirst {
		err = r.setTimeOut(time.Now().Add(serverConf.wsHeaderReadTimeout * time.Second))
		if err != nil {
			// error in setting timeout
			return &frame, newSetReadTimeoutError(err)
		}
	}

	numBytes, readBytes, err = r.read(1)
	if err != nil {
		return &frame, newReadError(err)
	}
	frame.secondByte = readBytes[0]

	if r._isClient && !frame.isMasked() {
		// return error, data is not masked by client
		return &frame, newExpectingMaskedFrameError(&frame)
	}
	if !r._isClient && frame.isMasked() {
		// return error, data is masked by server
		return &frame, newExpectingUnmaskedFrameError(&frame)
	}

	length := int(frame.secondByte&0x7f)
	if length == 0 && !frame.isControlFrame() {
		// return error, data frame length must be more than zero
		return &frame, newEmptyDataFrameError(&frame)
	}

	if length > 0x7d {

		if frame.isControlFrame() {
			// control frame must have length less than or equal to 125
			return &frame, newControlFrameLengthError(&frame)
		}

		if length == 0x7e {
			length = 2
		} else {
			length = 8
		}
		for {
			numBytes, readBytes, err = r.read(length)
			frame.lengthBytes = append(frame.lengthBytes, readBytes...)
			length -= numBytes
			if err != nil || length == 0 {
				break
			}
		}
		if err != nil {
			return &frame, newReadError(err)
		}
		length = frame.payloadLength()
		if frame.secondByte&0x7f == 0x7e && length < 126 {
			return &frame, newFramePayloadLengthError(&frame, true)
		} else if length <= 65536 {
			return &frame, newFramePayloadLengthError(&frame, false)
		}
	} else {
		if frame.opcode() == M_CLS && frame.payloadLength() == 1 {
			// return error, invaid data length received for closed frame
			return &frame, newCloseFrameLengthError(&frame)
		}
	}

	if frame.isMasked() {
		// read mask bytes
		length = 4
		for {
			numBytes, readBytes, err = r.read(length)
			frame.maskBytes = append(frame.maskBytes, readBytes...)
			length -= numBytes
			if err != nil || length == 0 {
				break
			}
		}
		if err != nil {
			return &frame, newReadError(err)
		}
	}

	return &frame, err
}

func (r *wsReader) readFrameData(frame *wsFrame) error {
	var(
		size int
		numBytes int
		cntBytes int
		buff []byte
		err error
	)
	size = frame.payloadLength()
	if size > 0 {

		sec := 10
		minBytes := sec * serverConf.wsMinByteRatePerSec

		err = r.setTimeOut(time.Now().Add(time.Duration(sec) * time.Second))
		if err != nil {
			// set read timeout failed
			return newSetReadTimeoutError(err)
		}

		cntBytes = 0

		for {
			numBytes, buff, err = r.read(size)

			cntBytes += numBytes

			if e, ok := err.(net.Error); ok && e.Timeout() {
				// timeout occured
				if cntBytes < minBytes {
					// return error, connection sent less data (numbytes bytes data) for 10 second
					// expecting minBytes (serverConf.wsMinByteRatePerSec per second)
					return newSlowDataReadError(cntBytes, sec)
				}
				err = r.setTimeOut(time.Now().Add(time.Duration(sec) * time.Second))
				if err != nil {
					// set read timeout failed
					return newSetReadTimeoutError(err)
				}
				err = nil
				cntBytes = 0
			}

			frame.data = append(frame.data, buff[:numBytes]...)
			size -= numBytes

			if err != nil || size == 0 {
				break
			}
		}
	}

	if err != nil {
		return newReadError(err)
	}
	return err
}
