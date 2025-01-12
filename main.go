package main

import (
	"fmt"
	"io"
	"net"
	"os"
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
				buf := make([]byte, 1024)

				// read message from client
				_, err = conn.Read(buf)
				if err != nil {
					if err == io.EOF {
						break
					}
					fmt.Println("error reading from client: ", err.Error())
					os.Exit(1)
				}

				// ignore request and send back a PONG
				conn.Write([]byte("+OKK\r\n"))
			}
		}(conn)
	}

}
