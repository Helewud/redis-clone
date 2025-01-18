package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/helewud/redis-clone/commands"
	"github.com/helewud/redis-clone/resp"
	"github.com/helewud/redis-clone/storage"
)

func main() {
	// Create a new server
	n, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer n.Close()

	store, err := restoreStoreBackup("storage.store")
	if err != nil {
		fmt.Println(err)
		return
	}

	// listen
	for {
		// Listen for connections
		conn, err := n.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handleConn(conn, store)
	}
}

func handleConn(conn net.Conn, store *storage.Aof) {
	defer conn.Close()

	for {
		value, err := validateRespInput(conn)
		if err != nil {
			fmt.Printf("resp input validation error: %q \n", err)
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, err := validateRespCommand(command)
		if err != nil {
			fmt.Printf("resp command error: %q \n", err)

			temp := resp.Value{T: resp.RespTString, String: ""}
			conn.Write(temp.Marshal())

			continue
		}

		result := handler(args)
		conn.Write(result.Marshal())

		if strings.Contains(command, "SET") {
			store.Write(*value)
		}
	}

}

func handleRespValue(value resp.Value) error {
	command := strings.ToUpper(value.Array[0].Bulk)
	args := value.Array[1:]

	handler, ok := commands.Handlers[command]
	if !ok {
		return fmt.Errorf("invalid command: %v", command)
	}

	handler(args)

	return nil
}

func restoreStoreBackup(backupPath string) (*storage.Aof, error) {
	store, err := storage.NewAof(backupPath)
	if err != nil {
		return nil, err
	}

	err = store.Read(handleRespValue)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func validateRespInput(conn net.Conn) (*resp.Value, error) {
	reader := resp.NewReader(conn)
	value, err := reader.Read()
	if err != nil {
		return nil, err
	}

	if value.T != resp.RespTArray {
		return nil, fmt.Errorf("invalid request, expected array")
	}

	if len(value.Array) == 0 {
		return nil, fmt.Errorf("invalid request, expected array length > 0")
	}

	return &value, nil
}

func validateRespCommand(command string) (commands.RespHandler, error) {
	handler, ok := commands.Handlers[command]
	if !ok {
		return nil, fmt.Errorf("invalid command: %v", command)
	}

	return handler, nil
}
