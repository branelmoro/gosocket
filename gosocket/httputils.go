package gosocket

import (
	"strings"
)


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