package main

import (
	"fmt"
	"net"

	"github.com/helewud/redis-clone/resp"
)

func main() {

	// Create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	// listen
	for {
		// Listen for connections
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go func(conn net.Conn) {
			defer conn.Close()

			for {
				reader := resp.NewReader(conn)
				value, err := reader.Read()
				if err != nil {
					fmt.Println(err)
					return
				}

				// print value to terminal
				fmt.Println(value)

				// ignore request and send back a OKK
				conn.Write([]byte("+OKK\r\n"))
			}
		}(conn)
	}

}
