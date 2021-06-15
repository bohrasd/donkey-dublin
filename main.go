package main

import (
	"fmt"
	"net"

	"github.com/donkey-dublin/http"
)

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

		go http.Process(conn)

	}

	// b := make([]byte, 1024)
	// in := os.Stdin
	// n, err := in.Read(b)

}
