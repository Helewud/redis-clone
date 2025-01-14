package resp

import (
	"bytes"
	"testing"
)

func TestIsEndOfLine(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  bool
	}{
		{
			name:  "end of line",
			input: []byte("\r\n"),
			want:  true,
		},
		{
			name:  "not end of line",
			input: []byte("\r"),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEndOfLine(tt.input); got != tt.want {
				t.Errorf("isEndOfLine(%q) = %v; want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestTrimEndOfLine(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{
			name:  "remove end of line",
			input: []byte("start\r\n"),
			want:  []byte("start"),
		},
		{
			name:  "return line as it is",
			input: []byte("start\r"),
			want:  []byte("start\r"),
		},
		{
			name:  "return simple text",
			input: []byte("start"),
			want:  []byte("start"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trimEndOfLine(tt.input)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("trimEndOfLine(%s) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestAppendEndOfLine(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{
			name:  "append end of line",
			input: []byte("start"),
			want:  []byte("start\r\n"),
		},
		{
			name:  "append end of line if  it exist",
			input: []byte("start\r\n"),
			want:  []byte("start\r\n\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := appendEndOfLine(tt.input)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("appendEndOfLine(%s) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
