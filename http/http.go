package http

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"os"
	"strconv"
	"strings"
)

var con *net.Conn

type Request struct {
	method, path, version string
	headers               map[string]string
}

func Process(conn net.Conn) {
	con = &conn

	// b := make([]byte, 1024)
	// n, err := conn.Read(b)
	bufrd := bufio.NewReader(conn)
	http_str, err := bufrd.ReadString('\n')
	if err != nil {
		fmt.Printf("cannot read: %v", err)
		return
	}

	firstLine := strings.Fields(http_str)
	req := Request{firstLine[0], firstLine[1], firstLine[2], make(map[string]string)}

	var header_str string
	var header_item []string
header:
	for {
		header_str, err = bufrd.ReadString('\n')
		header_item = strings.SplitN(header_str, ":", 2)
		// fmt.Printf("%#v\n", header_str)

		switch {
		case header_str == "\r\n":
			break header

		case len(header_item) == 2:
			req.headers[strings.TrimSpace(header_item[0])] = strings.TrimSpace(header_item[1])

		default:
			fmt.Println("Malformed header!")
			break header
		}
	}

	fmt.Printf("%#v\n", req)

	resp := Response{
		version: req.version,
		headers: map[string]string{
			"Server": "donkey-dublin/0.1",
		},
	}

	// word := string(b[:n])
	switch req.method {
	case "GET":
		resp.get(&req)
	case "POST":
	case "OPTIONS":
	case "PUT":
	case "DELETE":
	case "HEAD":
	default:
		resp.status = 400
		resp.body = strings.NewReader("Im sorry I just don't understand.")
		resp.render()
		log.Printf("unknown method %s", firstLine[0])
	}

	conn.Close()
}

type Response struct {
	version string
	status  int
	headers map[string]string
	body    io.Reader
}

var StatusMap = map[int]string{
	200: "OK",
	400: "Bad Request",
	404: "Not Found",
}

func (resp *Response) render() {

	bufwt := bufio.NewWriter(*con)
	firstLine := fmt.Sprintf("%s %d %s\r\n", resp.version, resp.status, StatusMap[resp.status])
	bufwt.Write([]byte(firstLine))

	//header
	for k, v := range resp.headers {
		bufwt.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
	}

	bufwt.Write([]byte("\r\n"))

	bufrd := bufio.NewReader(resp.body)
	io.Copy(bufwt, bufrd)
	bufwt.Write([]byte("\r\n"))

	bufwt.Flush()
}

func (resp *Response) get(req *Request) error {
	switch req.path {
	case "/test.txt":
		resp.status = 200

		file, err := os.Open("./www" + req.path)
		if err != nil {
			panic(err)
		}

		stat, err := file.Stat()
		if err != nil {
			panic(err)
		}

		resp.headers["Content-Length"] = strconv.FormatInt(stat.Size(), 10)
		resp.headers["Content-Type"] = mime.TypeByExtension(".txt")

		resp.body = file
		resp.render()

		file.Close()

	default:
		resp.status = 404
		resp.body = strings.NewReader("Sorry we don't have that file!")
		resp.render()
	}

	return nil
}
