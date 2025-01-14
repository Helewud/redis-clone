package resp

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestReader_readLine(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []byte
		wantErr bool
	}{
		{
			name:    "simple line",
			input:   "hello\r\n",
			want:    []byte("hello"),
			wantErr: false,
		},
		{
			name:    "empty line",
			input:   "\r\n",
			want:    []byte{},
			wantErr: false,
		},
		{
			name:    "incomplete line",
			input:   "hello",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(strings.NewReader(tt.input))
			got, _, err := r.readLine()

			if (err != nil) != tt.wantErr {
				t.Errorf("readLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !bytes.Equal(got, tt.want) {
					t.Errorf("readLine() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestReader_readInteger(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{
			name:    "positive number",
			input:   "42\r\n",
			want:    42,
			wantErr: false,
		},
		{
			name:    "negative number",
			input:   "-42\r\n",
			want:    -42,
			wantErr: false,
		},
		{
			name:    "zero",
			input:   "0\r\n",
			want:    0,
			wantErr: false,
		},
		{
			name:    "invalid number",
			input:   "abc\r\n",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(strings.NewReader(tt.input))
			got, err := r.parseInteger()

			if (err != nil) != tt.wantErr {
				t.Errorf("readInteger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got != tt.want {
					t.Errorf("readInteger() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestReader_readBulk(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Value
		wantErr bool
	}{
		{
			name:  "simple bulk string",
			input: "5\r\nhello\r\n",
			want: Value{
				T:    RespTBulk,
				Bulk: "hello",
			},
			wantErr: false,
		},
		{
			name:  "empty bulk string",
			input: "0\r\n\r\n",
			want: Value{
				T:    RespTBulk,
				Bulk: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(strings.NewReader(tt.input))
			got, err := r.readBulk()

			if (err != nil) != tt.wantErr {
				t.Errorf("readBulk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readBulk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_readArray(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Value
		wantErr bool
	}{
		{
			name:  "simple array",
			input: "2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			want: Value{
				T: RespTArray,
				Array: []Value{
					{T: RespTBulk, Bulk: "hello"},
					{T: RespTBulk, Bulk: "world"},
				},
			},
			wantErr: false,
		},
		{
			name:  "empty array",
			input: "0\r\n",
			want: Value{
				T: RespTArray,
			},
			wantErr: false,
		},
		{
			name:  "nested array",
			input: "2\r\n$5\r\nhello\r\n*1\r\n$5\r\nworld\r\n",
			want: Value{
				T: RespTArray,
				Array: []Value{
					{T: RespTBulk, Bulk: "hello"},
					{
						T: RespTArray,
						Array: []Value{
							{T: RespTBulk, Bulk: "world"},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(strings.NewReader(tt.input))
			got, err := r.readArray()

			if (err != nil) != tt.wantErr {
				t.Errorf("readArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_Read(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Value
		wantErr bool
	}{
		{
			name:  "read bulk string",
			input: string(RespBulk) + "5\r\nhello\r\n",
			want: Value{
				T:    RespTBulk,
				Bulk: "hello",
			},
			wantErr: false,
		},
		{
			name:  "read array",
			input: string(RespArray) + "2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			want: Value{
				T: RespTArray,
				Array: []Value{
					{T: RespTBulk, Bulk: "hello"},
					{T: RespTBulk, Bulk: "world"},
				},
			},
			wantErr: false,
		},
		{
			name:    "read invalid type",
			input:   "x5\r\nhello\r\n",
			want:    Value{},
			wantErr: false, // According to the implementation
		},
		{
			name:    "read empty input",
			input:   "",
			want:    Value{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(strings.NewReader(tt.input))
			got, err := r.Read()

			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestReader_EOF tests EOF handling
func TestReader_EOF(t *testing.T) {
	r := NewReader(strings.NewReader(""))
	_, err := r.Read()
	if err != io.EOF {
		t.Errorf("Expected EOF error, got %v", err)
	}
}
