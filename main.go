package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/helewud/redis-clone/commands"
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

				if value.T != resp.RespTArray {
					fmt.Println("Invalid request, expected array")
					continue
				}

				if len(value.Array) == 0 {
					fmt.Println("Invalid request, expected array length > 0")
					continue
				}

				command := strings.ToUpper(value.Array[0].Bulk)
				args := value.Array[1:]

				handler, ok := commands.Handlers[command]
				if !ok {
					fmt.Println("Invalid command: ", command)
					temp := resp.Value{T: resp.RespTString, String: ""}
					conn.Write(temp.Marshal())
					continue
				}

				result := handler(args)
				conn.Write(result.Marshal())
			}
		}(conn)
	}

}
