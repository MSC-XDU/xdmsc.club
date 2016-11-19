package signup

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"os"
	"path/filepath"
)

var infoDB *bolt.DB

var (
	profileBucket = []byte("profile")
	signUpBucket  = []byte("signup")
	// 指示时候是一个空白的信息
	EmptyInfoErr = errors.New("empty profile")
)

const (
	infoDBFile = resourcePath + "/info.db"
)

func init() {
	if err := os.MkdirAll(resourcePath, 0766); err != nil {
		panic(err)
	}

	dbPath, err := filepath.Abs(infoDBFile)
	db, err := bolt.Open(dbPath, 0766, nil)
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(profileBucket)
		_, err = tx.CreateBucketIfNotExists(signUpBucket)
		return err
	})

	if err != nil {
		panic(err)
	}

	infoDB = db

	go func() {
		<-cleanUp
		infoDB.Close()
	}()
}

// 获取此 token 下的个人信息，如果尚且为空，返回 Profile 的零值，并且返回 EmptyInfoErr 错误
func (token UserToken) GetProfile() (*Profile, error) {
	var p Profile
	return &p, token.getInfo(profileBucket, &p)
}

// 将此 token 下的个人信息，替换为提供的信息
func (token UserToken) PutProfile(profile *Profile) error {
	if !checkProfile(profile) {
		return InformationInvalidErr
	}
	return token.putInfo(profileBucket, profile)
}

// 获取此 token 下的报名信息，如果尚且为空，返回 SignUp 的零值，并且返回 EmptyInfoErr 错误
func (token UserToken) GetSignUp() (*SignUp, error) {
	var s SignUp
	return &s, token.getInfo(signUpBucket, &s)
}

// 将此 token 下的报名信息，替换为提供的信息
func (token UserToken) PutSignUp(signUp *SignUp) error {
	if !checkSignUp(signUp) {
		return InformationInvalidErr
	}
	return token.putInfo(signUpBucket, signUp)
}

func (token UserToken) getInfo(bucketName []byte, info interface{}) error {
	k := token.toBytes()
	err := infoDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		c := bucket.Cursor()
		k, v := c.Seek(k)
		if k == nil {
			return EmptyInfoErr
		}

		return json.Unmarshal(v, info)
	})
	return err
}

func (token UserToken) putInfo(bucketName []byte, info interface{}) error {
	k := token.toBytes()
	value, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return infoDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		return bucket.Put(k, value)
	})
}

func checkProfile(profile *Profile) bool {
	// TODO
	return true
}

func checkSignUp(signUp *SignUp) bool {
	// TODO
	return true
}

func (token UserToken) toBytes() []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(token))
	return buf
}
