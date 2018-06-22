package main

import (
	"fmt"
	"strings"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	// "reflect"
	// "unsafe"
)

// this is a comment

func main() {
    fmt.Println("Hello World")

    // byte2int()

    // var a error

    // a = nil

    // fmt.Println(a)
    // return

    // // readHTTP()
    // p,q,r := testy()
    // fmt.Println(p,q,r)

    // // readWebsocketFrame()

    // var (
    // 	b int
    // 	z uint
    // )

    // b = 200
    // z = (uint(b) << 45)
    // fmt.Println(z, b)

    // var abc int
    // abc = 2344342
    // // 100011 11000101 10010110
    // // 35 197 150
    // length_bytes := []byte{byte(abc >> 16),byte(abc >> 8),byte(abc >> 0)}
    // fmt.Println(length_bytes)
    // p := byte(abc >> 0)
    // fmt.Println(p)
    // p = byte(abc >> 8)
    // fmt.Println(p)
    // p = byte(abc >> 16)
    // fmt.Println(p)
    // var psc *int
    // fmt.Println(psc)

    // a := 123

    // psc = &a
    // fmt.Println(psc)

    // fmt.Println(reflect.TypeOf(psc))

    // var z *int
    // var add uintptr
    // add = 0xc4200120b8

    // z = (*int)(unsafe.Pointer(add))
    // fmt.Println(*z)

    decrept()

}

func decrept() {
	var a byte
	a  = 1
	fmt.Println(a, 0x1)
	switch (a) {
		case 0x0://continuation
		case 0x1://text
		case 0x2://binary
			fmt.Println("calling cb - ")
			break
		case 0x8://close
			break
		case 0x9://ping
			// send pong on ping
			break
		case 0xA://pong
			break
		default:
			fmt.Println("default case - opcode - ")
	}



	// a := []byte{129 146 115 227 5 160 27 134 105 204 28 195 99 210 28 142 37 194 1 140 114 211 22 145}
}


func byte2int() {
        aa := uint(0x7FFFFFFF)
        fmt.Println(aa)
        slice := []byte{0xFF, 0xFF, 0xFF, 0x7F}
        tt := binary.BigEndian.Uint32(slice)
        fmt.Println(tt)
}


func testy() (int,*int,int) {
	return 354,nil,767
}

func readWebsocketFrame() {
	
    // a := []byte{136,150,49,72,152,142,50,162,209,224,71,41,244,231,85,104,234,235,66,45,234,248,84,44,184,236,88,60}
    a := []byte{129,146,137,106,201,176,225,15,165,220,230,74,175,194,230,7,233,210,251,5,190,195,236,24}



    // if a[0]&0x8 >> 7 {
    // 	fmt.Println("true")
    // } else {
    // 	fmt.Println("false")
    // }


    fin := (a[0]&0x80 >> 7 == 1)
	if fin {
		fmt.Println("finished message")
	} else {
		fmt.Println("read next frame of message")
	}

    opCode, payloadType := opcode(a[0])

    fmt.Println(fin, opCode, payloadType)

	mask := a[1]&0x80 >> 7
	if mask == 1 {
		fmt.Println("valid mask from client...")
	} else {
		fmt.Println("mask not found... disconnect websocket..")
	}

	bytes_read := 2

	payloadLength := a[1]&0x7f
	// if payloadLength == 0x7e {
	if payloadLength == 126 {
		fmt.Println("next 16bit - 2 bytes(a[2]a[3]) is length but not supported")
		bytes_read += 2
	}
	// if payloadLength == 0x7f {
	if payloadLength == 127 {
		fmt.Println("next 64bit - 8 bytes(a[2]a[3]a[4]a[5]a[6]a[7]a[8]a[9]) is length but not supported")
		bytes_read += 8
	}

	var payload []byte

	fmt.Println("Payload Length is ", payloadLength)
	if mask == 1 {

		mask_key := a[bytes_read:bytes_read+4]
		fmt.Println("Mask key is ", mask_key)

		bytes_read += 4

		var ch byte
		for k,v := range a[bytes_read:] {
			ch = v^mask_key[k%4]
			payload = append(payload, ch)
		}

		fmt.Println("Payload is", string(payload), payload, a[bytes_read:])
	}



 //    fmt.Println(a[0]&0x80)
 //    b := a[0]&0x80 >> 7


}


func opcode(firstByte byte) (byte, string) {
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


func readHTTP() {


    a := []byte{71,69,84,32,47,101,99,104,111,32,72,84,84,80,47,49,46,49,13,10,72,111,115,116,58,32,108,111,99,97,108,104,111,115,116,58,51,51,51,51,13,10,85,115,101,114,45,65,103,101,110,116,58,32,77,111,122,105,108,108,97,47,53,46,48,32,40,88,49,49,59,32,85,98,117,110,116,117,59,32,76,105,110,117,120,32,120,56,54,95,54,52,59,32,114,118,58,54,48,46,48,41,32,71,101,99,107,111,47,50,48,49,48,48,49,48,49,32,70,105,114,101,102,111,120,47,54,48,46,48,13,10,65,99,99,101,112,116,58,32,116,101,120,116,47,104,116,109,108,44,97,112,112,108,105,99,97,116,105,111,110,47,120,104,116,109,108,43,120,109,108,44,97,112,112,108,105,99,97,116,105,111,110,47,120,109,108,59,113,61,48,46,57,44,42,47,42,59,113,61,48,46,56,13,10,65,99,99,101,112,116,45,76,97,110,103,117,97,103,101,58,32,101,110,45,71,66,44,101,110,59,113,61,48,46,53,13,10,65,99,99,101,112,116,45,69,110,99,111,100,105,110,103,58,32,103,122,105,112,44,32,100,101,102,108,97,116,101,13,10,83,101,99,45,87,101,98,83,111,99,107,101,116,45,86,101,114,115,105,111,110,58,32,49,51,13,10,79,114,105,103,105,110,58,32,104,116,116,112,58,47,47,108,111,99,97,108,104,111,115,116,13,10,83,101,99,45,87,101,98,83,111,99,107,101,116,45,69,120,116,101,110,115,105,111,110,115,58,32,112,101,114,109,101,115,115,97,103,101,45,100,101,102,108,97,116,101,13,10,83,101,99,45,87,101,98,83,111,99,107,101,116,45,75,101,121,58,32,56,116,52,56,67,98,117,89,118,83,66,104,106,48,86,104,98,89,111,83,53,81,61,61,13,10,67,111,110,110,101,99,116,105,111,110,58,32,107,101,101,112,45,97,108,105,118,101,44,32,85,112,103,114,97,100,101,13,10,80,114,97,103,109,97,58,32,110,111,45,99,97,99,104,101,13,10,67,97,99,104,101,45,67,111,110,116,114,111,108,58,32,110,111,45,99,97,99,104,101,13,10,85,112,103,114,97,100,101,58,32,119,101,98,115,111,99,107,101,116,13,10,13,10}

    // req, headers, body, err := parseRequest(&a)

    // fmt.Println(req, headers, body, err)


    fmt.Println(string(a))
    req, req_len, err := getRequest(&a)

    fmt.Println(req, req_len, err)
    fmt.Println(req, req_len, err)
    fmt.Println(req, req_len, err)
    fmt.Println(req, req_len, err)
    fmt.Println(req, req_len, err)
    fmt.Println(err)
    fmt.Println(err)
    fmt.Println(req)

	if err == "" {
		headers, err :=getHeaders(&a,(req_len+1))
	    fmt.Println(headers, err)
	    fmt.Println(headers, err)
	    fmt.Println(headers, err)

		if v1, v2 := headers["Upgrade"], headers["Sec-WebSocket-Key"]; v1 != "" && v1 == "websocket" && v2 != "" && req[0] == "GET" {
			fmt.Println("ebsocket")
		}
	}

	// if val, ok :=headers["Upgrade"]; ok && val == "websocket" && _, ok := headers["Sec-WebSocket-Key"]; ok {
	// 	fmt.Println("ebsocket")
	// }

	// if v1, v2 := headers["Upgrade"], headers["Sec-WebSocket-Key"]; v1 != nil && v1 == "websocket" && v2 != nil {
	// 	fmt.Println("ebsocket")
	// }


    // str := "dGhlIHNhbXBsZSBub25jZQ=="
    // str1 := getSecWebSocketAccept(str)
    // fmt.Println(str)
    // fmt.Println(str1)
	
}


func getRequest(rb *[]byte) ([]string, int, string) {
	var (
		req []string
		err string
		prev_byte byte
		req_len int
	)
	// err = nil
	// req = make([]string)

	read_bytes := *rb

	prev_byte = 0
	for k, v := range read_bytes {
		if prev_byte == 13 && v == 10 {
			req = strings.Split(string(read_bytes[:k]), " ")
			req_len = k
			break
		}
		prev_byte = v
	}

	if len(req) > 0 {
		switch(req[0]) {
			case "GET":
			case "POST":
			case "PUT":
			case "DELETE":
				break
			default:
				err = "Invalid HTTP websocket request method" 
				break
		}
	}

	return req, req_len, err
}

func getHeaders(rb *[]byte, req_len int) (map[string]string, string) {

	var err string

	read_bytes := *rb

	header_lines := strings.Split(string(read_bytes[(req_len+1):]), "\r\n")

	headers := make(map[string]string)

	for _, element := range header_lines {
		hd := strings.SplitN(element, ":", 2)
		if len(hd) == 2 {
			headers[strings.Trim(hd[0], " ")] = strings.Trim(hd[1], " ")
		}
	}
	return headers, err
}


func getSecWebSocketAccept(s string) string {
	str := append([]byte(s), []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")...)
    // s := "dGhlIHNhbXBsZSBub25jZQ==258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
    h := sha1.New()
    h.Write(str)
    bs := h.Sum(nil)
    return base64.StdEncoding.EncodeToString(bs)
}



func parseRequest(read_bytes *[]byte) ([]string, map[string]string, []byte, string) {

	var(
		req []string
		headers map[string]string
		err string
		body []byte
	)

	req_data := strings.SplitN(string(*read_bytes), "\r\n", 2)

	req = strings.Split(req_data[0], " ")

	if len(req) != 3 {
		err = "Invalid HTTP request"
		return req, headers, body, err
	} else {

		switch(req[0]) {
			case "GET":
				req_data = strings.SplitN(req_data[1], "\r\n\r\n", 2)
				headers, err = parseHTTPHeader(&req_data[0])
				// for k,v := range headers {
				// 	fmt.Println(k,"-----",v)
				// }
				break
			case "POST":
			case "PUT":
			case "DELETE":
			default:
				err = "Invalid HTTP websocket request method" 
				break
		}
	}

	return req, headers, body, err
}




func parseHTTPHeader(header_data *string) (map[string]string, string) {

	var err string

	header_lines := strings.Split(*header_data, "\r\n")

	headers := make(map[string]string)

	for _, element := range header_lines {
		hd := strings.SplitN(element, ":", 2)
		if len(hd) == 2 {
			headers[strings.Trim(hd[0], " ")] = strings.Trim(hd[1], " ")
		}
	}

	return headers, err
}