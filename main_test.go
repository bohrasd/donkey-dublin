package main

import (
	"net"
	"testing"
)

func TestHttp(t *testing.T) {
	// go main()

	tests := []struct {
		req, resp string
	}{
		{"GET a", "abc"},
		{"GET c", "ghi"},
		{"ALL", "key: c val: ghi\nkey: a val: asd\nkey: b val: def"},
	}

	resp_b := make([]byte, 1024)
	var resp string
	for _, test := range tests {
		conn, _ := net.Dial("tcp4", "127.0.0.1:30000")

		conn.Write([]byte(test.req))

		n, _ := conn.Read(resp_b)
		resp = string(resp_b[:n])
		if resp != test.resp {
			t.Errorf("%v", resp)
		}
		conn.Close()
	}
}
