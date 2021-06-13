package main

import (
	"fmt"
	"net"
	"strings"
)

func process(conn net.Conn) {

	// dict to search
	dic := map[string]string{
		"a": "asd",
		"b": "def",
		"c": "ghi",
	}

	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		fmt.Printf("cannot read: %v", err)
		return
	}

	word := string(b[:n])
	fmt.Printf("%q\n", word)

	body := strings.SplitN(word, " ", 2)
	switch body[0] {
	case "GET":
		stripped_str := strings.TrimSpace(body[1])
		ans, ok := dic[stripped_str]
		if !ok {
			conn.Write([]byte("ERROR undefined"))
		} else {
			conn.Write([]byte("ANSWER " + ans))
		}

	case "SET":
		item := strings.SplitN(body[1], " ", 2)
		dic[item[0]] = item[1]
	case "CLEAR":
		dic = make(map[string]string)

	case "ALL":
		for k, v := range dic {
			conn.Write([]byte(fmt.Sprintf("key: %s val: %v\n", k, v)))
		}
	}

	conn.Close()

	return
}

func main() {

	listener, err := net.Listen("tcp4", "127.0.0.1:30000")
	if err != nil {
		fmt.Printf("cannot listen: %v", err)
		return
	}

	for {

		conn, err := listener.Accept()

		if err != nil {
			fmt.Printf("cannot accept: %v", err)
			return
		}

		go process(conn)

	}

	// b := make([]byte, 1024)
	// in := os.Stdin
	// n, err := in.Read(b)

}
