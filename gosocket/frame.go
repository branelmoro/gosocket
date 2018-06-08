package gosocket

func getOpCode(firstByte byte) (byte, string) {
	var payloadType string
	opcode := firstByte & 0x0f
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
			payloadType = "reserved for non-control"
	}
	return opcode, payloadType
}

func (c *Conn) readMessage() *[]byte {

	var message []byte

	for {
		read_bytes, fin, err := c.conn.readFrame()

		message = append(message, read_bytes...)

		if fin {
			break
		}
	}
	return &message
}


func (c *Conn)readFrame() *[]byte, bool {

	var(
		fin bool
		rsv1 bool
		rsv2 bool
		rsv3 bool
		mask bool
		opcode byte
		payloadLength byte
		mask_key [4]byte
		read_bytes []byte
		buff *[]byte
	)

	fin, rsv1, rsv2, rsv3, opcode = c.readFirstByteFromFrame()

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
			panic("Invalid frame opcode...reserved for non-control")
	}

	mask, payloadLength := c.readSecondByteFromFrame()

	if mask {
		if payloadLength == 126 {
			panic("next 16bit - 2 bytes(a[2]a[3]) is length but not supported")
		}
		if payloadLength == 127 {
			panic("next 64bit - 8 bytes(a[2]a[3]a[4]a[5]a[6]a[7]a[8]a[9]) is length but not supported")
		}

		_, buff = c.readBytes(4)
		mask_key = *buff

		var (
			num_bytes int
			cnt_byte int
		)
		cnt_byte = 0

		for num_bytes, buff = range c.readBytes(payloadLength) {

			for _, ch := range *buff {
				read_bytes = append(read_bytes, (ch^mask_key[cnt_byte%4]))
				cnt_byte += 1
			}

			if num_bytes == payloadLength {
				break
			} else {
				payloadLength -= num_bytes
			}
		}

		return &read_bytes, fin
	} else {
		panic("mask not found... disconnect websocket..")
	}
}


func (c *Conn)readFirstByteFromFrame() (bool,bool,bool,bool,byte) {
	_, buff := c.readBytes(1)
	read_bytes := *buff
	fin     := (read_bytes[0]&0x80 >> 7 == 1)
	rsv1    := (read_bytes[0]&0x40 >> 6 == 1)
	rsv2    := (read_bytes[0]&0x20 >> 5 == 1)
	rsv3    := (read_bytes[0]&0x10 >> 4 == 1)
	opcode  := read_bytes[0]&0x0f >> 7
	return fin,rsv1,rsv2,rsv3,opcode
}

func (c *Conn)readSecondByteFromFrame() (bool,byte) {
	_, buff := c.readBytes(1)
	read_bytes := *buff
	mask := (read_bytes[0]&0x80 >> 7 == 1)
	payloadLength := read_bytes[0]&0x7f
	return mask, payloadLength
}


func (c *Conn)readBytes(buff_size int) int, *[]byte {

	buff := make([]byte, buff_size)
	num_bytes, err := c.conn.Read(buff)

	if err != nil {
		if err != io.EOF {
			panic(err)
			// fmt.Println("read error:", err)
		}
	}

	return num_bytes, &buff
}





func (c *Conn) Read() *[]byte {

	var read_bytes []byte

	buff_size := 1

	timeoutDuration := 1 * time.Millisecond
	fmt.Println("Time-----", timeoutDuration)
	c.conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	buff := make([]byte, buff_size)
	for {
		num_bytes, err := c.conn.Read(buff)
		// fmt.Println("Bytes received:", num_bytes, err, string(buff), time.Now())

		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}

		read_bytes = append(read_bytes, buff...)

		if num_bytes < buff_size {
			break
		} else {
			c.conn.SetReadDeadline(time.Now().Add(timeoutDuration))
		}
	}
	return &read_bytes
}