package gosocket

import(
	"time"
	"fmt"
	"io"
)

func (c *Conn) readMessages() {
	var(
		message *[]byte
		byteCnt int
		err error
		msg_len int
	)
	for{
		message, msg_len, byteCnt, err = c.readMessage()
		if err != nil {
			if byteCnt > 0 {
				c.Close()
			}
			break
		} else {
			OnMessage(c, message)
		}
	}
	fmt.Println(msg_len)
}

func (c *Conn) readMessage() (*[]byte, int, int, error) {

	var(
		fin bool
		frame_payload *[]byte
		payloadLength int
		num_bytes int
		message []byte
		is_first bool
		byteCnt int
		err error
		msg_len int
	)

	is_first = true
	byteCnt = 0
	msg_len = 0

	fmt.Println(is_first)

	for {
		fin, frame_payload, payloadLength, num_bytes, err = c.readFrame()

		byteCnt += num_bytes
		msg_len += payloadLength

		if err != nil {
			break
		}

		message = append(message, *frame_payload...)

		if fin {
			break
		}
	}

	return &message, msg_len, byteCnt, err
}


func (c *Conn)readFrame() (bool, *[]byte, int, int, error) {

	var(
		fin bool
		rsv1 bool
		rsv2 bool
		rsv3 bool
		opcode byte
		mask bool
		payloadLength int
		mask_key []byte
		frame_payload []byte
		buff *[]byte
		err error
		byteCnt int
		num_bytes int
		payloadType string
	)

	byteCnt = 0

	fin, rsv1, rsv2, rsv3, opcode, num_bytes, err = c.readFirstByteFromFrame()

	fmt.Println(rsv3,rsv2,rsv1,payloadType)
	byteCnt += num_bytes
	if err != nil {
		return fin, &frame_payload, payloadLength, byteCnt, err
	}

	switch (opcode) {
		case 0x0:
			payloadType = "continuation"
			break
		case 0x1:
			payloadType = "text"
			break
		case 0x2:
			payloadType = "binary"
			break
		case 0x8:
			payloadType = "connection close"
			// 136 130 245 134 144 67 246 110
			// 10001000 10000010 11110101 10000110 10010000 01000011 11110110 01101110
			// 11110101 10000110
			// 11110110 01101110
			// .................
			// 00000011 11101000

			// [136 128 132 166 99 33]
			// 10001000 10000000 10000100 10100110 01100011 00100001
			break
		case 0x9:
			payloadType = "ping"
			break
		case 0xA:
			payloadType = "pong"
			break
		default:
			err = NewWsError(PAYLOAD_LENGTH_ERROR, "Invalid frame opcode...reserved for non-control")
			return fin, &frame_payload, payloadLength, byteCnt, err
	}

	mask, payloadLength, num_bytes, err = c.readSecondByteFromFrame()
	byteCnt += num_bytes
	if err != nil {
		return fin, &frame_payload, payloadLength, byteCnt, err
	}


	if mask {

		c.conn.SetReadDeadline(time.Time{})

		fmt.Println("initial payloadLength - ", payloadLength)

		if payloadLength == 126 {
			num_bytes, buff, err = c.readBytes(2)
			byteCnt += num_bytes
			if err != nil {
				return fin, &frame_payload, payloadLength, byteCnt, err
			}

			len_bytes := *buff
			payloadLength = ((int(len_bytes[0]) << 8) | int(len_bytes[1]))
			if payloadLength < 126 {
				err = NewWsError(PAYLOAD_LENGTH_ERROR, "Invalid payload length in 16 bit")
				return fin, &frame_payload, payloadLength, byteCnt, err
			}

			fmt.Println("----------here----------")

			// err = NewWsError(PAYLOAD_LENGTH_ERROR, "next 16bit - 2 bytes(a[2]a[3]) is length but not supported")
			// return fin, &frame_payload, payloadLength, byteCnt, err
		} else if payloadLength == 127 {

			num_bytes, buff, err = c.readBytes(8)
			byteCnt += num_bytes
			if err != nil {
				return fin, &frame_payload, payloadLength, byteCnt, err
			}

			len_bytes := *buff
			payloadLength = ((int(len_bytes[0]) << 56) | (int(len_bytes[1]) << 48) | (int(len_bytes[2]) << 40) | (int(len_bytes[3]) << 32) | (int(len_bytes[4]) << 24) | (int(len_bytes[5]) << 16) | (int(len_bytes[6]) << 8) | int(len_bytes[7]))

			fmt.Println("64 bit Payload length - ", payloadLength)
			if payloadLength < 65535 {
				err = NewWsError(PAYLOAD_LENGTH_ERROR, "Invalid payload length in 64 bit")
				return fin, &frame_payload, payloadLength, byteCnt, err
			}

			// err = NewWsError(PAYLOAD_LENGTH_ERROR, "next 64bit - 8 bytes(a[2]a[3]a[4]a[5]a[6]a[7]a[8]a[9]) is length but not supported")
			// return fin, &frame_payload, payloadLength, byteCnt, err
		}

		fmt.Println("final payloadLength - ", payloadLength)

		num_bytes, buff, err = c.readBytes(4)
		byteCnt += num_bytes
		if err != nil {
			return fin, &frame_payload, payloadLength, byteCnt, err
		}
		mask_key = *buff

		fmt.Println(mask_key)

		if payloadLength > 0 {

			var (
				mask_count int
			)
			mask_count = 0

			num_bytes, buff, err = c.readPayload(payloadLength)
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

		return fin, &frame_payload, payloadLength, byteCnt, err
	} else {
		panic("mask not found... disconnect websocket..")
	}
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
		opcode  = read_bytes[0]&0x0f >> 7
	}
	return fin, rsv1, rsv2, rsv3, opcode, num_bytes, err
}

func (c *Conn)readSecondByteFromFrame() (bool, int, int, error) {
	var(
		mask bool
		payloadLength int
		err error
		buff *[]byte
		num_bytes int
	)

	c.conn.SetReadDeadline(time.Now().Add(time.Second))
	num_bytes, buff, err = c.readBytes(1)

	if err == nil {
		read_bytes      := *buff
		mask            = (read_bytes[0]&0x80 >> 7 == 1)
		payloadLength   = int(read_bytes[0]&0x7f)
	}
	return mask, payloadLength, num_bytes, err
}

func (c *Conn)readPayload(size int) (int, *[]byte, error) {
	var(
		err error
		num_bytes int
		cntbytes int
		read_bytes []byte
		buff *[]byte
	)
	cntbytes = 0
	for{
		num_bytes, buff, err = c.readBytes(size)
		if err != nil {
			break
		}

		read_bytes = append(read_bytes, *buff...)

		cntbytes += num_bytes

		size -= num_bytes

		if size == 0 {
			break
		}
	}
	return cntbytes, &read_bytes, err
}

func (c *Conn)readBytes(size int) (int, *[]byte, error) {
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


