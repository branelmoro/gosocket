package main

import (
	"fmt"
	"strings"
	"crypto/sha1"
	"encoding/base64"
)

// this is a comment

func main() {
    fmt.Println("Hello World")


    a := []byte{71,69,84,32,47,101,99,104,111,32,72,84,84,80,47,49,46,49,13,10,72,111,115,116,58,32,108,111,99,97,108,104,111,115,116,58,51,51,51,51,13,10,85,115,101,114,45,65,103,101,110,116,58,32,77,111,122,105,108,108,97,47,53,46,48,32,40,88,49,49,59,32,85,98,117,110,116,117,59,32,76,105,110,117,120,32,120,56,54,95,54,52,59,32,114,118,58,54,48,46,48,41,32,71,101,99,107,111,47,50,48,49,48,48,49,48,49,32,70,105,114,101,102,111,120,47,54,48,46,48,13,10,65,99,99,101,112,116,58,32,116,101,120,116,47,104,116,109,108,44,97,112,112,108,105,99,97,116,105,111,110,47,120,104,116,109,108,43,120,109,108,44,97,112,112,108,105,99,97,116,105,111,110,47,120,109,108,59,113,61,48,46,57,44,42,47,42,59,113,61,48,46,56,13,10,65,99,99,101,112,116,45,76,97,110,103,117,97,103,101,58,32,101,110,45,71,66,44,101,110,59,113,61,48,46,53,13,10,65,99,99,101,112,116,45,69,110,99,111,100,105,110,103,58,32,103,122,105,112,44,32,100,101,102,108,97,116,101,13,10,83,101,99,45,87,101,98,83,111,99,107,101,116,45,86,101,114,115,105,111,110,58,32,49,51,13,10,79,114,105,103,105,110,58,32,104,116,116,112,58,47,47,108,111,99,97,108,104,111,115,116,13,10,83,101,99,45,87,101,98,83,111,99,107,101,116,45,69,120,116,101,110,115,105,111,110,115,58,32,112,101,114,109,101,115,115,97,103,101,45,100,101,102,108,97,116,101,13,10,83,101,99,45,87,101,98,83,111,99,107,101,116,45,75,101,121,58,32,56,116,52,56,67,98,117,89,118,83,66,104,106,48,86,104,98,89,111,83,53,81,61,61,13,10,67,111,110,110,101,99,116,105,111,110,58,32,107,101,101,112,45,97,108,105,118,101,44,32,85,112,103,114,97,100,101,13,10,80,114,97,103,109,97,58,32,110,111,45,99,97,99,104,101,13,10,67,97,99,104,101,45,67,111,110,116,114,111,108,58,32,110,111,45,99,97,99,104,101,13,10,85,112,103,114,97,100,101,58,32,119,101,98,115,111,99,107,101,116,13,10,13,10}

    req, headers, body, err := parseRequest(&a)

    fmt.Println(req, headers, body, err)



    str := "dGhlIHNhbXBsZSBub25jZQ=="
    str1 := getSecWebSocketAccept(str)
    fmt.Println(str)
    fmt.Println(str1)

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