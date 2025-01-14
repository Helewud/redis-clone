package commands

import (
	"sync"

	"github.com/helewud/redis-clone/resp"
)

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{
			T:      resp.RespTError,
			String: "ERR wrong number of arguments for 'HSET' command",
		}
	}

	rkey := args[0].Bulk
	pkey := args[1].Bulk
	value := args[2].Bulk

	res := resp.Value{
		T:      resp.RespTString,
		String: "OK",
	}

	HSETsMu.Lock()
	_, ok := HSETs[rkey]
	if !ok {
		HSETs[rkey] = map[string]string{
			pkey: value,
		}
	}
	if ok {
		HSETs[rkey][pkey] = value
	}
	HSETsMu.Unlock()

	return res
}

func hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{
			T:      resp.RespTError,
			String: "ERR wrong number of arguments for 'HGET' command",
		}
	}

	rkey := args[0].Bulk
	pkey := args[1].Bulk

	res := resp.Value{T: resp.RespTNull}

	HSETsMu.RLock()
	if value, ok := HSETs[rkey][pkey]; ok {
		res.T = resp.RespTBulk
		res.Bulk = value
	}
	HSETsMu.RUnlock()

	return res
}

func hgetall(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{
			T:      resp.RespTError,
			String: "ERR wrong number of arguments for 'HGETALL' command",
		}
	}

	rkey := args[0].Bulk

	val := []resp.Value{}
	res := resp.Value{T: resp.RespTNull}

	HSETsMu.RLock()
	if value, ok := HSETs[rkey]; ok {
		for k, v := range value {
			val = append(val, resp.Value{T: resp.RespTBulk, Bulk: k})
			val = append(val, resp.Value{T: resp.RespTBulk, Bulk: v})
		}
		res.T = resp.RespTArray
		res.Array = val
	}
	HSETsMu.RUnlock()

	return res
}
