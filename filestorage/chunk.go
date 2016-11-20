package filestorage

import (
	"io"
	"os"
)

type Chunk struct {
	Offset   int64
	Len      int64
	At       string
	FileName string
}

func (c Chunk) GetBytes() ([]byte, error) {
	f, err := os.Open(c.At)
	if err != nil {
		return nil, err
	}
	b := make([]byte, c.Len)
	_, err = f.ReadAt(b, c.Offset)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return b, nil
}
