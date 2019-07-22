package gosocket


func applyMask(data []byte, mask []byte) {
	var mk byte = 0x0
	for i, ch := range data {
		mk &= 0x3
		data[i] = ch^mask[mk]
		mk += 1
	}
}

type wsFrame struct{
	firstByte byte
	secondByte byte
	lengthBytes []byte
	maskBytes []byte
	data []byte
}

func (frame *wsFrame) getType() string {
	switch (frame.opcode()) {
		case M_TXT:
			return "text"
		case M_BIN:
			return "binary"
		case M_CLS:
			return "close"
		case M_PING:
			return "ping"
		case M_PONG:
			return "pong"
		default:
			return "unidentified"
	}
}

func (frame *wsFrame) fin() bool {
	return frame.firstByte&0x80 == 0x80
}

func (frame *wsFrame) rsv1() bool {
	return frame.firstByte&0x40 == 0x40
}

func (frame *wsFrame) rsv2() bool {
	return frame.firstByte&0x20 == 0x20
}

func (frame *wsFrame) rsv3() bool {
	return frame.firstByte&0x10 == 0x10
}

func (frame *wsFrame) opcode() byte {
	return frame.firstByte&0x0f
}

func (frame *wsFrame) isControlFrame() bool {
	return frame.firstByte&0x0f == M_CLS ||
		frame.firstByte&0x0f == M_PING ||
		frame.firstByte&0x0f == M_PONG
}

func (frame *wsFrame) isMasked() bool {
	return frame.secondByte&0x80 == 0x80
}

func (frame *wsFrame) payloadLength() int {
	payloadLength := int(frame.secondByte&0x7f)
	if payloadLength > 0x7d {
		shift := uint8(8 * (len(frame.lengthBytes)-1))
		for _, ch := range frame.lengthBytes {
			payloadLength |= (int(ch) << shift)
			shift -= 8
		}
		return payloadLength
	} else {
		return payloadLength
	}
}

func (frame *wsFrame) length() int {
	return 2 + len(frame.lengthBytes) + len(frame.maskBytes) +  frame.payloadLength()
}

func (frame *wsFrame) payload() []byte {
	return frame.data
}

func (frame *wsFrame) unMask(data []byte) {
	applyMask(data[len(data)-frame.payloadLength():], frame.maskBytes)
}

func (frame *wsFrame) toBytes() []byte {
	var data []byte
	data = append(data, frame.firstByte, frame.secondByte)
	data = append(data, frame.lengthBytes...)
	data = append(data, frame.maskBytes...)
	return append(data, frame.data...)
}


func generateMask() []byte {
	return []byte{34,74,89,13}
}


type wFrame struct{
	fin bool
	rsv1 bool
	rsv2 bool
	rsv3 bool
	opcode byte
	isMasked bool
	data []byte
}

func (frame *wFrame) toBytes() []byte {
	var(
		rawBytes []byte
		tempByte byte
		length int
	)

	// set first byte
	tempByte = frame.opcode
	if frame.fin { tempByte |= 0x80 }
	if frame.rsv1 { tempByte |= 0x40 }
	if frame.rsv2 { tempByte |= 0x20 }
	if frame.rsv3 { tempByte |= 0x10 }
	rawBytes = append(rawBytes, tempByte)

	// set masking bit
	if frame.isMasked {
		tempByte = 0x80
	} else {
		tempByte = 0x00
	}

	length = len(frame.data)
	
	// set length
	if length <= 125 {
		tempByte |= byte(length)
		rawBytes = append(rawBytes, tempByte)
	} else {
		if length <= 65536 {
			tempByte |= 0x7e
			rawBytes = append(rawBytes, tempByte, byte(length >> 8), byte(length))
		} else {
			tempByte |= 0x7f
			rawBytes = append(rawBytes, tempByte, byte(length >> 56), byte(length >> 48), byte(length >> 40), byte(length >> 32), byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length))
		}
	}

	// set data in frame
	if frame.isMasked {
		maskBytes := generateMask()
		// set mask bytes
		rawBytes = append(rawBytes, maskBytes...)
		length = len(rawBytes)
		rawBytes = append(rawBytes, frame.data...)
		applyMask(rawBytes[length:], maskBytes)
	} else {
		rawBytes = append(rawBytes, frame.data...)
	}
	return rawBytes
}

