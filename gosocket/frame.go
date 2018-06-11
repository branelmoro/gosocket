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

			err = NewWsError(PAYLOAD_LENGTH_ERROR, "next 16bit - 2 bytes(a[2]a[3]) is length but not supported")
			return fin, &frame_payload, payloadLength, byteCnt, err
		}
		if payloadLength == 127 {

			num_bytes, buff, err = c.readBytes(8)
			byteCnt += num_bytes
			if err != nil {
				return fin, &frame_payload, payloadLength, byteCnt, err
			}

			len_bytes := *buff
			payloadLength = ((int(len_bytes[0]) << 56) | (int(len_bytes[0]) << 48) | (int(len_bytes[0]) << 40) | (int(len_bytes[0]) << 32) | (int(len_bytes[0]) << 24) | (int(len_bytes[0]) << 16) | (int(len_bytes[0]) << 8) | int(len_bytes[1]))
			if payloadLength < 65535 {
				err = NewWsError(PAYLOAD_LENGTH_ERROR, "Invalid payload length in 64 bit")
				return fin, &frame_payload, payloadLength, byteCnt, err
			}

			err = NewWsError(PAYLOAD_LENGTH_ERROR, "next 64bit - 8 bytes(a[2]a[3]a[4]a[5]a[6]a[7]a[8]a[9]) is length but not supported")
			return fin, &frame_payload, payloadLength, byteCnt, err
		}

		num_bytes, buff, err = c.readBytes(4)
		byteCnt += num_bytes
		if err != nil {
			return fin, &frame_payload, payloadLength, byteCnt, err
		}
		mask_key = *buff

		fmt.Println(mask_key)

		var (
			mask_count int
		)
		mask_count = 0

		cntBytesToRead := payloadLength

		for {
			num_bytes, buff, err = c.readBytes(cntBytesToRead) 

			byteCnt += num_bytes

			if err != nil {
				break
			}

			for _, ch := range *buff {
				fmt.Println(mask_count)
				frame_payload = append(frame_payload, (ch^mask_key[mask_count]))
				mask_count += 1
				if mask_count == 4 {
					mask_count = 0
				}
			}

			if num_bytes == cntBytesToRead {
				break
			} else {
				cntBytesToRead -= num_bytes
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
	num_bytes, buff, err = c.readBytes(1)

	if err == nil {
		read_bytes      := *buff
		mask            = (read_bytes[0]&0x80 >> 7 == 1)
		payloadLength   = int(read_bytes[0]&0x7f)
	}
	return mask, payloadLength, num_bytes, err
}

func (c *Conn)readBytes(buff_size int) (int, *[]byte, error) {
	var(
		err error
		num_bytes int
	)

	timeoutDuration := time.Duration(buff_size) * time.Millisecond
	c.conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	buff := make([]byte, buff_size)
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


