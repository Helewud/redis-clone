package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Reader struct {
	reader *bufio.Reader
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(rd)}
}

func (r *Reader) readLine() (line []byte, charCount int, err error) {
	for {
		char, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		charCount++

		line = append(line, char)

		if isEndOfLine(line) {
			return trimEndOfLine(line), charCount, nil
		}
	}
}

func (r *Reader) parseInteger() (value int, err error) {
	line, _, err := r.readLine()
	if err != nil {
		return 0, err
	}

	int64Val, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, err
	}

	return int(int64Val), nil
}

func (r *Reader) readArray() (Value, error) {
	v := Value{t: RespTArray}

	// get array length
	len, err := r.parseInteger()
	if err != nil {
		return v, err
	}

	// foreach line, parse and read the value
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		// append parsed value to array
		v.array = append(v.array, val)
	}

	return v, nil
}

func (r *Reader) readBulk() (Value, error) {
	v := Value{
		t: RespTBulk,
	}

	// get bulk char length
	len, err := r.parseInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	// remove CRLF from line
	r.readLine()

	return v, nil
}

func (r *Reader) Read() (Value, error) {
	symbol, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch Symbol(symbol) {
	case RespArray:
		return r.readArray()
	case RespBulk:
		return r.readBulk()
	default:
		fmt.Printf("Unknown resp type symbol: %v", string(symbol))
		return Value{}, nil
	}
}
