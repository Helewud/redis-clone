package resp

type Symbol string

const (
	RespString  Symbol = "+"
	RespBulk    Symbol = "$"
	RespArray   Symbol = "*"
	RespError   Symbol = "-"
	RespInteger Symbol = ":"
	RespNnull   Symbol = "$-1\r\n"
)

type Type string

const (
	RespTArray  Type = "array"
	RespTBulk   Type = "bulk"
	RespTNull   Type = "null"
	RespTError  Type = "error"
	RespTString Type = "string"
)

type Value struct {
	T      Type
	String string
	Number int
	Bulk   string
	Array  []Value
}
