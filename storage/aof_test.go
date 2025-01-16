package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/helewud/redis-clone/resp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewAof tests the creation of new AOF instance
func TestNewAof(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.aof")

	t.Run("successful creation", func(t *testing.T) {
		aof, err := NewAof(path)
		require.NoError(t, err)
		require.NotNil(t, aof)
		require.NotNil(t, aof.file)
		require.NotNil(t, aof.reader)

		err = aof.Close()
		require.NoError(t, err)
	})

	t.Run("invalid path", func(t *testing.T) {
		invalidPath := filepath.Join(tmpDir, "nonexistent", "test.aof")
		aof, err := NewAof(invalidPath)
		require.Error(t, err)
		require.Nil(t, aof)
	})
}

// TestAofWrite tests the Write operation
func TestAofWrite(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "write_test.aof")

	aof, err := NewAof(path)
	require.NoError(t, err)
	defer aof.Close()

	t.Run("single write", func(t *testing.T) {
		value := resp.Value{
			T:    resp.RespTBulk,
			Bulk: "test data",
		}

		err := aof.Write(value)
		require.NoError(t, err)

		// Verify file contents
		data, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Contains(t, string(data), "test data")
	})

	t.Run("concurrent writes", func(t *testing.T) {
		var wg sync.WaitGroup
		numWrites := 20

		for i := 0; i < numWrites; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				value := resp.Value{
					T:    resp.RespTBulk,
					Bulk: fmt.Sprintf("data-%d", i),
				}
				err := aof.Write(value)
				require.NoError(t, err)
			}(i)
		}

		wg.Wait()

		// Verify all writes were successful
		data, err := os.ReadFile(path)
		require.NoError(t, err)
		content := string(data)

		for i := 0; i < numWrites; i++ {
			assert.Contains(t, content, fmt.Sprintf("data-%d", i))
		}
	})
}

// TestAofRead tests the Read operation
func TestAofRead(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "read_test.aof")

	aof, err := NewAof(path)
	require.NoError(t, err)
	defer aof.Close()

	// Write test data
	testValues := []resp.Value{
		{T: resp.RespTBulk, Bulk: "test1"},
		{T: resp.RespTBulk, Bulk: "test2"},
		{T: resp.RespTBulk, Bulk: "test3"},
	}

	for _, value := range testValues {
		err := aof.Write(value)
		require.NoError(t, err)
	}

	t.Run("read all values", func(t *testing.T) {
		var readValues []string
		err := aof.Read(func(value resp.Value) error {
			readValues = append(readValues, value.Bulk)
			return nil
		})
		require.NoError(t, err)

		assert.Equal(t, len(testValues), len(readValues))
		for i, value := range testValues {
			assert.Equal(t, value.Bulk, readValues[i])
		}
	})

	t.Run("read with error callback", func(t *testing.T) {
		expectedErr := fmt.Errorf("test error")
		err := aof.Read(func(value resp.Value) error {
			return expectedErr
		})
		assert.Equal(t, expectedErr, err)
	})
}

// TestAofSync tests the automatic sync functionality
func TestAofSync(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "sync_test.aof")

	aof, err := NewAof(path)
	require.NoError(t, err)
	defer aof.Close()

	// Write test data
	value := resp.Value{T: resp.RespTBulk, Bulk: "sync test"}
	err = aof.Write(value)
	require.NoError(t, err)

	// Wait for sync
	time.Sleep(2 * time.Second)

	// Verify data was synced
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}

// TestAofClose tests the Close operation
func TestAofClose(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "close_test.aof")

	aof, err := NewAof(path)
	require.NoError(t, err)

	// Write some data
	value := resp.Value{T: resp.RespTBulk, Bulk: "close test"}
	err = aof.Write(value)
	require.NoError(t, err)

	// Close the file
	err = aof.Close()
	require.NoError(t, err)

	// Verify file is closed by attempting to write
	err = aof.Write(value)
	assert.Error(t, err)
}

// TestAofLocking tests the mutex locking behavior
func TestAofLocking(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "lock_test.aof")

	aof, err := NewAof(path)
	require.NoError(t, err)
	defer aof.Close()

	t.Run("concurrent operations", func(t *testing.T) {
		var wg sync.WaitGroup
		numOps := 100

		// Concurrent writes
		for i := 0; i < numOps; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				value := resp.Value{
					T:    resp.RespTBulk,
					Bulk: fmt.Sprintf("concurrent-%d", i),
				}
				err := aof.Write(value)
				require.NoError(t, err)
			}(i)
		}

		// Concurrent reads
		for i := 0; i < numOps; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := aof.Read(func(value resp.Value) error {
					return nil
				})
				require.NoError(t, err)
			}()
		}

		wg.Wait()
	})
}

// TestAofRecovery tests the recovery of data after restart
func TestAofRecovery(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "recovery_test.aof")

	// Write initial data
	aof1, err := NewAof(path)
	require.NoError(t, err)

	initialData := []resp.Value{
		{T: resp.RespTBulk, Bulk: "data1"},
		{T: resp.RespTBulk, Bulk: "data2"},
	}

	for _, value := range initialData {
		err = aof1.Write(value)
		require.NoError(t, err)
	}

	err = aof1.Close()
	require.NoError(t, err)

	// Reopen file and verify data
	aof2, err := NewAof(path)
	require.NoError(t, err)
	defer aof2.Close()

	var recoveredData []string
	err = aof2.Read(func(value resp.Value) error {
		recoveredData = append(recoveredData, value.Bulk)
		return nil
	})
	require.NoError(t, err)

	assert.Equal(t, len(initialData), len(recoveredData))
	for i, value := range initialData {
		assert.Equal(t, value.Bulk, recoveredData[i])
	}
}
