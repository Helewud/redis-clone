package commands

import (
	"sync"

	"github.com/helewud/redis-clone/resp"
)

func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{
			T:      resp.RespTString,
			String: "PONG",
		}
	}

	return resp.Value{
		T:      resp.RespTString,
		String: args[0].Bulk,
	}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{
			T:      resp.RespTError,
			String: "ERR wrong number of arguments for 'SET' command",
		}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return resp.Value{
		T:      resp.RespTString,
		String: "OK",
	}
}

func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{
			T:      resp.RespTError,
			String: "ERR wrong number of arguments for 'GET' command",
		}
	}

	key := args[0].Bulk
	res := resp.Value{T: resp.RespTNull}

	SETsMu.RLock()
	if value, ok := SETs[key]; ok {
		res.T = resp.RespTBulk
		res.Bulk = value
	}
	SETsMu.RUnlock()

	return res
}
