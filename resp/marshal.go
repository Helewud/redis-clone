package resp

import (
	"bytes"
	"strconv"
)

func (v Value) marshalString() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(string(RespString))
	buffer.WriteString(v.str)
	return appendEndOfLine(buffer.Bytes())
}

func (v Value) marshalBulk() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(string(RespBulk))
	buffer.WriteString(strconv.Itoa(len(v.bulk)))
	buffer.WriteString("\r\n")
	buffer.WriteString(v.bulk)
	return appendEndOfLine(buffer.Bytes())
}

func (v Value) marshalArray() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(string(RespArray))
	buffer.WriteString(strconv.Itoa(len(v.array)))
	buffer.WriteString("\r\n")
	for _, item := range v.array {
		buffer.Write(item.Marshal())
	}
	return buffer.Bytes()
}

func (v Value) marshalError() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(string(RespError))
	buffer.WriteString(v.str)
	return appendEndOfLine(buffer.Bytes())
}

func (v Value) marshalNull() []byte {
	return []byte(RespNnull)
}

func (v Value) Marshal() []byte {
	switch v.t {
	case RespTArray:
		return v.marshalArray()
	case RespTBulk:
		return v.marshalBulk()
	case RespTNull:
		return v.marshalNull()
	case RespTError:
		return v.marshalError()
	case RespTString:
		return v.marshalString()
	default:
		return []byte{}
	}
}
