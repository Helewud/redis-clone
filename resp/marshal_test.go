package resp

import (
	"bytes"
	"testing"
)

func TestValue_Marshal(t *testing.T) {
	tests := []struct {
		name string
		v    Value
		want []byte
	}{
		{
			name: "simple string",
			v: Value{
				t:   RespTString,
				str: "hello",
			},
			want: []byte("+hello\r\n"),
		},
		{
			name: "bulk string",
			v: Value{
				t:    RespTBulk,
				bulk: "hello",
			},
			want: []byte("$5\r\nhello\r\n"),
		},
		{
			name: "empty bulk string",
			v: Value{
				t:    RespTBulk,
				bulk: "",
			},
			want: []byte("$0\r\n\r\n"),
		},
		{
			name: "null value",
			v: Value{
				t: RespTNull,
			},
			want: []byte("$-1\r\n"),
		},
		{
			name: "error message",
			v: Value{
				t:   RespTError,
				str: "Error occurred",
			},
			want: []byte("-Error occurred\r\n"),
		},
		{
			name: "simple array",
			v: Value{
				t: RespTArray,
				array: []Value{
					{t: RespTString, str: "hello"},
					{t: RespTString, str: "world"},
				},
			},
			want: []byte("*2\r\n+hello\r\n+world\r\n"),
		},
		{
			name: "empty array",
			v: Value{
				t:     RespTArray,
				array: []Value{},
			},
			want: []byte("*0\r\n"),
		},
		{
			name: "nested array",
			v: Value{
				t: RespTArray,
				array: []Value{
					{t: RespTString, str: "hello"},
					{
						t: RespTArray,
						array: []Value{
							{t: RespTString, str: "world"},
						},
					},
				},
			},
			want: []byte("*2\r\n+hello\r\n*1\r\n+world\r\n"),
		},
		{
			name: "mixed array",
			v: Value{
				t: RespTArray,
				array: []Value{
					{t: RespTString, str: "hello"},
					{t: RespTBulk, bulk: "world"},
					{t: RespTNull},
					{t: RespTError, str: "test error"},
				},
			},
			want: []byte("*4\r\n+hello\r\n$5\r\nworld\r\n$-1\r\n-test error\r\n"),
		},
		{
			name: "unknown type",
			v: Value{
				t: "unknown",
			},
			want: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v.Marshal()
			if !bytes.Equal(got, tt.want) {
				t.Errorf("Marshal() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValue_marshalString(t *testing.T) {
	tests := []struct {
		name string
		v    Value
		want []byte
	}{
		{
			name: "simple string",
			v:    Value{str: "hello"},
			want: []byte("+hello\r\n"),
		},
		{
			name: "empty string",
			v:    Value{str: ""},
			want: []byte("+\r\n"),
		},
		{
			name: "string with spaces",
			v:    Value{str: "hello world"},
			want: []byte("+hello world\r\n"),
		},
		{
			name: "string with special chars",
			v:    Value{str: "hello\nworld"},
			want: []byte("+hello\nworld\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.t = "string"
			got := tt.v.marshalString()
			if !bytes.Equal(got, tt.want) {
				t.Errorf("marshalString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValue_marshalBulk(t *testing.T) {
	tests := []struct {
		name string
		v    Value
		want []byte
	}{
		{
			name: "simple bulk",
			v:    Value{bulk: "hello"},
			want: []byte("$5\r\nhello\r\n"),
		},
		{
			name: "empty bulk",
			v:    Value{bulk: ""},
			want: []byte("$0\r\n\r\n"),
		},
		{
			name: "bulk with spaces",
			v:    Value{bulk: "hello world"},
			want: []byte("$11\r\nhello world\r\n"),
		},
		{
			name: "bulk with special chars",
			v:    Value{bulk: "hello\nworld"},
			want: []byte("$11\r\nhello\nworld\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.t = "bulk"
			got := tt.v.marshalBulk()
			if !bytes.Equal(got, tt.want) {
				t.Errorf("marshalBulk() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValue_marshalError(t *testing.T) {
	tests := []struct {
		name string
		v    Value
		want []byte
	}{
		{
			name: "simple error",
			v:    Value{str: "Error occurred"},
			want: []byte("-Error occurred\r\n"),
		},
		{
			name: "empty error",
			v:    Value{str: ""},
			want: []byte("-\r\n"),
		},
		{
			name: "error with special chars",
			v:    Value{str: "Error: invalid\ncharacter"},
			want: []byte("-Error: invalid\ncharacter\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.t = "error"
			got := tt.v.marshalError()
			if !bytes.Equal(got, tt.want) {
				t.Errorf("marshalError() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValue_marshalNull(t *testing.T) {
	v := Value{t: "null"}
	want := []byte("$-1\r\n")
	got := v.marshalNull()
	if !bytes.Equal(got, want) {
		t.Errorf("marshalNull() = %q, want %q", got, want)
	}
}

// Benchmark tests for performance analysis
func BenchmarkValue_Marshal(b *testing.B) {
	v := Value{
		t: RespTArray,
		array: []Value{
			{t: RespTString, str: "hello"},
			{t: RespTBulk, bulk: "world"},
			{t: RespTNull},
			{t: RespTError, str: "test error"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Marshal()
	}
}
