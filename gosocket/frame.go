package gosocket

import(
	"time"
	"fmt"
	"io"
)

type Message struct {
	opcode byte
	data *[]byte
}

func (c *Message) GetData() *[]byte {
	return c.data
}

func (c *Conn) WriteMessage(data *[]byte) {
	c.write(true, false, false, false, 0x01, data)
}

func (c *Conn) close() {
	p := *c.poller
	p.Stop(c.desc)
	c.desc.Close()
	c.conn.Close()
}

func (c *Conn) write(fin bool, rsv1 bool, rsv2 bool, rsv3 bool, opcode byte, data *[]byte) {
	var (
		message []byte
		length_bytes []byte
		length int
		// opcode byte
		first_byte byte
		second_byte byte
	)
	if data == nil {
		length = 0
	} else {
		length = len(*data)
	}

	first_byte = opcode

	if fin {
		first_byte |= 0x80
	}
	if rsv1 {
		first_byte |= 0x40
	}
	if rsv2 {
		first_byte |= 0x20
	}
	if rsv3 {
		first_byte |= 0x10
	}

	if length <= 125 {
		second_byte = byte(length)
	} else {
		if length <= 65535 {
			second_byte = byte(126)
			length_bytes = []byte{byte(length >> 8),byte(length >> 0)}
		} else {
			second_byte = byte(127)
			length_bytes = []byte{byte(length >> 48),byte(length >> 48),byte(length >> 40),byte(length >> 32),byte(length >> 24),byte(length >> 16),byte(length >> 8),byte(length)}
		}
	}

	message = append(message, first_byte, second_byte)
	message = append(message, length_bytes...)
	message = append(message, *data...)
	c.conn.Write(message)
}

func (c *Conn) CloseWebsocket(code int) {
	length_bytes := []byte{byte(code >> 8),byte(code >> 0)}
	c.write(true, false, false, false, 0x08, &length_bytes)
	c.close()
}

func (c *Conn) Ping(data *[]byte) {
	c.write(true, false, false, false, 0x09, data)
}

func (c *Conn) Pong(data *[]byte) {
	c.write(true, false, false, false, 0x0A, data)
}

func (c *Conn) readMessages() {
	var(
		message *Message
		byteCnt int
		err error
		msg_len int
	)
	for{
		message, msg_len, byteCnt, err = c.readMessage()
		if err != nil {
			if byteCnt > 0 {
				c.CloseWebsocket(PROTOCOL_ERROR)
			}
			break
		} else {
			fmt.Println("here-- - opcode - ", message.opcode, 0x1)
			switch (message.opcode) {
				case 0x0://continuation
					break
				case 0x1://text
					OnMessage(c, message)
					break
				case 0x2://binary
					OnMessage(c, message)
					break
				case 0x8://close
					OnClose(c, message)
					c.CloseWebsocket(NORMAL_CLOSURE)
					break
				case 0x9://ping
					// send pong on ping
					msg := *message
					c.Pong(msg.data)
					break
				case 0xA://pong
					break
				default:
					fmt.Println("default case - opcode - ", message.opcode)
			}
		}
		fmt.Println(msg_len)
	}
}

func (c *Conn) readMessage() (*Message, int, int, error) {

	var(
		fin bool
		frame_payload *[]byte
		payloadLength int
		num_bytes int
		msg_bytes []byte
		is_first_frame bool
		byteCnt int
		err error
		msg_len int
		opcode byte
		msg_opcode byte
	)

	is_first_frame = true
	byteCnt = 0
	msg_len = 0

	for {
		fin, frame_payload, opcode, payloadLength, num_bytes, err = c.readFrame()

		byteCnt += num_bytes
		msg_len += payloadLength

		if err != nil {
			break
		}

		if is_first_frame {
			is_first_frame = false
			msg_opcode = opcode
		} else if opcode != 0x0 {
			err = NewWsError(OPCODE_ERROR, "Invalid opcode in continuation frame...")
			break
		}

		msg_bytes = append(msg_bytes, *frame_payload...)

		if fin {
			break
		}
	}

	return &Message{opcode:msg_opcode,data:&msg_bytes}, msg_len, byteCnt, err
}


func (c *Conn)readFrame() (bool, *[]byte, byte, int, int, error) {

	var(
		fin bool
		rsv1 bool
		rsv2 bool
		rsv3 bool
		opcode byte
		second_byte byte
		mask bool
		payloadLength int
		mask_key []byte
		frame_payload []byte
		buff *[]byte
		err error
		byteCnt int
		num_bytes int
	)

	byteCnt = 0

	fin, rsv1, rsv2, rsv3, opcode, num_bytes, err = c.readFirstByteFromFrame()

	fmt.Println(rsv3,rsv2,rsv1)
	byteCnt += num_bytes
	if err != nil {
		return fin, &frame_payload, opcode, payloadLength, byteCnt, err
	}

	fmt.Println("Opcode is - ", opcode)

	switch (opcode) {
		case 0x0:
			break
		case 0x1:
			break
		case 0x2:
			break
		case 0x8:
			// 136 130 245 134 144 67 246 110
			// 10001000 10000010 11110101 10000110 10010000 01000011 11110110 01101110
			// 11110101 10000110
			// 11110110 01101110
			// .................
			// 00000011 11101000

			// [136 128 132 166 99 33]
			// 10001000 10000000 10000100 10100110 01100011 00100001
		case 0x9:
			break
		case 0xA:
			break
		default:
			err = NewWsError(OPCODE_ERROR, "Invalid opcode in frame...reserved for non-control")
			return fin, &frame_payload, opcode, payloadLength, byteCnt, err
	}


	num_bytes, buff, err = c.readBytes(1)
	byteCnt += num_bytes
	if err != nil {
		return fin, &frame_payload, opcode, payloadLength, byteCnt, err
	}
	read_bytes      := *buff
	second_byte     = read_bytes[0]
	mask            = (second_byte&0x80 >> 7 == 1)

	num_bytes, payloadLength, err = c.readPayloadLength(second_byte)
	byteCnt += num_bytes
	if err != nil {
		return fin, &frame_payload, opcode, payloadLength, byteCnt, err
	}

	if mask {

		c.conn.SetReadDeadline(time.Time{})

		fmt.Println("final payloadLength - ", payloadLength)

		num_bytes, buff, err = c.readBytes(4)
		byteCnt += num_bytes
		if err != nil {
			return fin, &frame_payload, opcode, payloadLength, byteCnt, err
		}
		mask_key = *buff

		fmt.Println(mask_key)

		if payloadLength > 0 {

			var (
				mask_count int
			)
			mask_count = 0

			num_bytes, buff, err = c.readBytes(payloadLength)
			fmt.Println(num_bytes, err)

			byteCnt += num_bytes
			if err == nil {
				for _, ch := range *buff {
					frame_payload = append(frame_payload, (ch^mask_key[mask_count]))
					mask_count += 1
					if mask_count == 4 {
						mask_count = 0
					}
				}
			}
		}
	} else {
		err = NewWsError(MASK_BIT_ERROR, "Mask bit not found")
	}
	return fin, &frame_payload, opcode, payloadLength, byteCnt, err
}


func (c *Conn)readFirstByteFromFrame() (bool, bool, bool, bool, byte, int, error) {
	var(
		fin bool
		rsv1 bool
		rsv2 bool
		rsv3 bool
		opcode byte
		err error
		buff *[]byte
		num_bytes int
	)

	c.conn.SetReadDeadline(time.Now().Add(time.Second))
	num_bytes, buff, err = c.readBytes(1)
	if err == nil {
		read_bytes := *buff
		fin     = (read_bytes[0]&0x80 >> 7 == 1)
		rsv1    = (read_bytes[0]&0x40 >> 6 == 1)
		rsv2    = (read_bytes[0]&0x20 >> 5 == 1)
		rsv3    = (read_bytes[0]&0x10 >> 4 == 1)
		opcode  = read_bytes[0]&0x0f
	}
	return fin, rsv1, rsv2, rsv3, opcode, num_bytes, err
}

func (c *Conn)readPayloadLength(second_byte byte) (int, int, error) {

	var(
		err error
		num_bytes int
		buff *[]byte
		shift uint8
		payloadLength int
		ch byte
	)

	payloadLength   = int(second_byte&0x7f)
	num_bytes       = 0

	if payloadLength > 125 {
		if payloadLength == 126 {
			num_bytes = 2
			shift = 8
		} else if payloadLength == 127 {
			num_bytes = 8
			shift = 56
		}
		
		num_bytes, buff, err = c.readBytes(num_bytes)

		if err != nil {
			return num_bytes, payloadLength, err
		}

		payloadLength = 0
		for _, ch = range *buff {
			payloadLength |= (int(ch) << shift)
			shift -= 8
		}

		if num_bytes == 2 && payloadLength < 126 {
			err = NewWsError(PAYLOAD_LENGTH_ERROR, "Invalid payload length in 16 bit")
			return num_bytes, payloadLength, err
		}
		if num_bytes == 8 && payloadLength < 65535 {
			err = NewWsError(PAYLOAD_LENGTH_ERROR, "Invalid payload length in 64 bit")
			return num_bytes, payloadLength, err
		}
	}

	return num_bytes, payloadLength, err
}

func (c *Conn)readBytes(size int) (int, *[]byte, error) {
	var(
		err error
		num_bytes int
		cntbytes int
		read_bytes []byte
		buff *[]byte
	)
	cntbytes = 0
	byte_size := size
	for{
		num_bytes, buff, err = c.readBytesInBuffer(size)
		if err != nil {
			break
		}

		cntbytes += num_bytes

		if num_bytes == byte_size {
			return cntbytes, buff, err
		}

		bytes_read := *buff

		read_bytes = append(read_bytes, bytes_read[:(num_bytes-1)]...)

		size -= num_bytes

		if size == 0 {
			break
		}
	}
	return cntbytes, &read_bytes, err
}

func (c *Conn)readBytesInBuffer(size int) (int, *[]byte, error) {
	var(
		err error
		num_bytes int
	)
	buff := make([]byte, size)

	num_bytes, err = c.conn.Read(buff)

	if err != nil {
		if err != io.EOF {
			// panic(err)
			fmt.Println("read error:", err)
			err = NewWsError(READ_ERROR, err.Error())
		} else {
			err = NewWsError(EOF_ERROR, err.Error())
		}
	}

	return num_bytes, &buff, err
}


