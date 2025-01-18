package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/helewud/redis-clone/resp"
	"github.com/helewud/redis-clone/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleRespValue(t *testing.T) {
	tests := []struct {
		name        string
		value       resp.Value
		expectError bool
	}{
		{
			name: "valid SET command",
			value: resp.Value{
				T: resp.RespTArray,
				Array: []resp.Value{
					{T: resp.RespTBulk, Bulk: "SET"},
					{T: resp.RespTBulk, Bulk: "key"},
					{T: resp.RespTBulk, Bulk: "value"},
				},
			},
			expectError: false,
		},
		{
			name: "invalid command",
			value: resp.Value{
				T: resp.RespTArray,
				Array: []resp.Value{
					{T: resp.RespTBulk, Bulk: "INVALID"},
					{T: resp.RespTBulk, Bulk: "key"},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleRespValue(tt.value)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRestoreStoreBackup(t *testing.T) {
	// Setup test file
	tmpDir := t.TempDir()
	backupFile := "storage.store"
	backupFilePath := filepath.Join(tmpDir, backupFile)
	defer os.Remove(backupFile)

	t.Run("successful restore", func(t *testing.T) {
		// Create and populate test store
		originalStore, err := storage.NewAof(backupFilePath)
		require.NoError(t, err)

		testValues := []resp.Value{
			{
				T: resp.RespTArray,
				Array: []resp.Value{
					{T: resp.RespTBulk, Bulk: "SET"},
					{T: resp.RespTBulk, Bulk: "key1"},
					{T: resp.RespTBulk, Bulk: "value1"},
				},
			},
			{
				T: resp.RespTArray,
				Array: []resp.Value{
					{T: resp.RespTBulk, Bulk: "SET"},
					{T: resp.RespTBulk, Bulk: "key2"},
					{T: resp.RespTBulk, Bulk: "value2"},
				},
			},
		}

		for _, value := range testValues {
			err = originalStore.Write(value)
			require.NoError(t, err)
		}
		require.NoError(t, originalStore.Close())

		// Test restore using the same file
		restoredStore, err := restoreStoreBackup(backupFilePath)
		require.NoError(t, err)
		require.NotNil(t, restoredStore)
		defer restoredStore.Close()

		// Verify that the restored store contains the expected data
		for _, testValue := range testValues {
			key := testValue.Array[1].Bulk
			expectedValue := testValue.Array[2].Bulk

			getHandler, err := validateRespCommand("GET")
			require.NoError(t, err)

			getResult := getHandler([]resp.Value{{T: resp.RespTBulk, Bulk: key}})
			assert.Equal(t, resp.Value{T: resp.RespTBulk, Bulk: expectedValue}, getResult,
				fmt.Sprintf("Restored value for key %s does not match", key))
		}
	})
}

func TestValidateRespInput(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectError bool
		expectValue *resp.Value
	}{
		{
			name:        "valid SET command",
			input:       []byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"),
			expectError: false,
			expectValue: &resp.Value{
				T: resp.RespTArray,
				Array: []resp.Value{
					{T: resp.RespTBulk, Bulk: "SET"},
					{T: resp.RespTBulk, Bulk: "key"},
					{T: resp.RespTBulk, Bulk: "value"},
				},
			},
		},
		{
			name:        "empty array",
			input:       []byte("*0\r\n"),
			expectError: true,
		},
		{
			name:        "invalid type",
			input:       []byte("+Simple String\r\n"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientConn, serverConn := net.Pipe()
			defer clientConn.Close()
			defer serverConn.Close()

			// Write input in a separate goroutine
			go func(data []byte) {
				_, _ = clientConn.Write(data)
				_ = clientConn.Close() // signal EOF to the other end
			}(tt.input)

			// Now read from `serverConn` as if it's the "server side"
			value, err := validateRespInput(serverConn)

			if tt.expectError {
				assert.Error(t, err, "expected an error but got none")
				return
			}

			require.NoError(t, err, "unexpected error when reading input")
			assert.Equal(t, tt.expectValue, value, "parsed Value mismatch")
		})
	}
}

// Test validateRespCommand
func TestValidateRespCommand(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		expectError bool
	}{
		{
			name:        "valid SET command",
			command:     "SET",
			expectError: false,
		},
		{
			name:        "valid GET command",
			command:     "GET",
			expectError: false,
		},
		{
			name:        "invalid command",
			command:     "INVALID",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, err := validateRespCommand(tt.command)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, handler)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, handler)
			}
		})
	}
}
