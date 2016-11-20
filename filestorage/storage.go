package filestorage

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/satori/go.uuid"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

const (
	chunkFileSize = 10 * 1048576
	resourcePath  = "data/xdmsc.club/files"
	dbFile        = resourcePath + "/files.db"
)

var fileStorage *Storage
var cleanUp = make(chan os.Signal)

func init() {
	if err := os.MkdirAll(resourcePath, 0766); err != nil {
		panic(err)
	}

	dbPath, err := filepath.Abs(dbFile)
	db, err := bolt.Open(dbPath, 0766, nil)
	if err != nil {
		panic(err)
	}

	fileBucket := []byte("file")
	chunkBucket := []byte("chunk")
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(fileBucket)
		_, err = tx.CreateBucketIfNotExists(chunkBucket)
		return err
	})

	if err != nil {
		panic(err)
	}

	fileStorage = &Storage{
		DB:            db,
		ChunkBucket:   chunkBucket,
		FileBucket:    fileBucket,
		ResourceRoot:  resourcePath,
		ChunkFileSize: chunkFileSize,
		FileMode:      0766,
	}

	go func() {
		<-cleanUp
		fileStorage.DB.Close()
	}()

	signal.Notify(cleanUp, os.Interrupt, os.Kill)
}

type Storage struct {
	DB            *bolt.DB
	ChunkBucket   []byte
	FileBucket    []byte
	ResourceRoot  string
	ChunkFileSize int64
	FileMode      os.FileMode
}

func GetChunk(key string) (Chunk, error) {
	return fileStorage.GetChunk(key)
}

func SaveFile(name string, value []byte) (string, error) {
	return fileStorage.SaveFile(name, value)
}

func (s *Storage) GetChunk(key string) (Chunk, error) {
	k := []byte(key)
	var chunk Chunk
	err := s.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.ChunkBucket)
		value := bucket.Get(k)
		if value == nil {
			return errors.New("file not found")
		}
		return json.Unmarshal(value, &chunk)
	})
	chunk.At, _ = filepath.Abs(strings.Join([]string{s.ResourceRoot, chunk.At}, "/"))
	return chunk, err
}

func (s *Storage) SaveFile(name string, value []byte) (string, error) {
	hash := md5.Sum(value)
	size := int64(len(value))
	chunk := Chunk{
		Len:      int64(size),
		FileName: name,
	}

	err := s.DB.Update(func(tx *bolt.Tx) error {
		fileBucket := tx.Bucket(s.FileBucket)
		chunkBucket := tx.Bucket(s.ChunkBucket)
		chunk.At, chunk.Offset = findChunkFile(fileBucket.Cursor(), size, s.ChunkFileSize)
		chunkJSON, err := json.Marshal(chunk)
		if err != nil {
			return err
		}

		path, _ := filepath.Abs(strings.Join([]string{s.ResourceRoot, chunk.At}, "/"))
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, s.FileMode)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = file.WriteAt(value, chunk.Offset)
		if err != nil {
			return err
		}

		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(chunk.Offset+chunk.Len))
		err = fileBucket.Put([]byte(chunk.At), b)
		if err != nil {
			return err
		}

		err = chunkBucket.Put(hash[:], chunkJSON)
		if err != nil {
			return err
		}
		chunk.At = path
		return nil
	})

	return string(hash[:]), err
}

func findChunkFile(c *bolt.Cursor, size, max int64) (string, int64) {
	if size >= max {
		return uuid.NewV4().String(), 0
	}

	for k, v := c.First(); k != nil; k, v = c.Next() {
		l := int64(binary.BigEndian.Uint64(v))
		if size <= max-l {
			return string(k), l
		}
	}

	return uuid.NewV4().String(), 0
}

// 用于清理数据库相关的操作
func CleanUp() {
	close(cleanUp)
}
