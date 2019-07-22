package gosocket

type HttpResponse struct {
	code string
	Headers map[string]string
}

func (r *HttpResponse) toBytes() []byte {
	return []byte{}
}
