package commands

import "github.com/helewud/redis-clone/resp"

type RespHandler = func([]resp.Value) resp.Value

var Handlers = map[string]RespHandler{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"DEL":     del,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}
