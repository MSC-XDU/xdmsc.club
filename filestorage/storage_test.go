package filestorage

import (
	"os"
	"testing"
)

const (
	testFileName    = "test.txt"
	testFileContent = "Hello!"
)

var fileKey string

func TestSaveFile(t *testing.T) {
	bytes := []byte(testFileContent)
	var err error
	fileKey, err = SaveFile(testFileName, bytes)
	if err != nil {
		t.Error(err)
	}
}

func TestGetChunk(t *testing.T) {
	chunk, err := GetChunk(fileKey)
	if err != nil {
		t.Fatal(err)
	}
	if chunk.FileName != testFileName {
		t.Error("name not match")
	}
	if chunk.Len != int64(len([]byte(testFileContent))) {
		t.Error("len not match")
	}
	if chunk.Offset != 0 {
		t.Error("offset not match")
	}
}

func TestChunk_GetBytes(t *testing.T) {
	chunk, err := GetChunk(fileKey)
	if err != nil {
		t.Error(err)
	}
	bytes, err := chunk.GetBytes()
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != testFileContent {
		t.Errorf("bad content %s\n", string(bytes))
	}
}

func TestMain(m *testing.M) {
	m.Run()
	CleanUp()
	os.RemoveAll("data")
}
