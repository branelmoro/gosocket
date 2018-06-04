package gosocket


func getSecWebSocketAccept(s string) string {
	str := append([]byte(s), []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")...)
    // s := "dGhlIHNhbXBsZSBub25jZQ==258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
    h := sha1.New()
    h.Write(str)
    bs := h.Sum(nil)
    return base64.StdEncoding.EncodeToString(bs)
}