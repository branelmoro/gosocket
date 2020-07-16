package gosocket

func isWhiteSpace(ch byte) bool {
	return ch == 0x20 || ch == 0x9
}

func isControlChar(ch byte) bool {
	return (ch < 0x20 && ch != 0x9) || ch == 0x7f
}