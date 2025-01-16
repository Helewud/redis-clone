package storage

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"

	"github.com/helewud/redis-clone/resp"
)

type Aof struct {
	file   *os.File
	reader *bufio.Reader
	mu     sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file:   f,
		reader: bufio.NewReader(f),
	}

	// Start a goroutine to sync AOF to disk every 1 second
	go func() {
		for {
			aof.mu.Lock()

			aof.file.Sync()

			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

func (aof *Aof) Write(value resp.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (aof *Aof) Read(callback func(value resp.Value) error) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	// Reset file pointer to beginning of file
	_, err := aof.file.Seek(0, 0)
	if err != nil {
		return err
	}

	// Reset the reader with the file
	aof.reader = bufio.NewReader(aof.file)
	reader := resp.NewReader(aof.reader)

	for {
		value, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := callback(value); err != nil {
			return err
		}
	}

	return nil
}
