package signup

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"os"
	"path/filepath"
)

var loginDB *bolt.DB
var loginBucket = []byte("login")

const (
	loginBDFile = resourcePath + "/login.db"
	PasswordLen = 6
)

func init() {
	if err := os.MkdirAll(resourcePath, 0766); err != nil {
		panic(err)
	}

	dbPath, err := filepath.Abs(loginBDFile)
	db, err := bolt.Open(dbPath, 0766, nil)
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(loginBucket)
		return err
	})

	if err != nil {
		panic(err)
	}

	loginDB = db

	go func() {
		<-cleanUp
		loginDB.Close()
	}()
}

func checkPassword(passwd string) error {
	if len(passwd) < PasswordLen {
		return PasswordTooShortErr
	}
	return nil
}

func (info Register) Register() (UserToken, error) {
	if err := checkPassword(info.Password); err != nil {
		return 0, err
	}
	uname := []byte(info.Username)
	value, err := json.Marshal(info)
	if err != nil {
		return 0, err
	}
	err = loginDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(loginBucket))

		e := bucket.Get(uname)
		if e != nil {
			return UserExistErr
		}

		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}
		info.Id = UserToken(id)

		return bucket.Put(uname, value)
	})
	return info.Id, err
}

func (info Register) LogIn() (UserToken, error) {
	var id UserToken = 0
	uname := []byte(info.Username)
	err := loginDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(loginBucket)
		value := bucket.Get(uname)
		if value == nil {
			return UserNotExistErr
		}

		var data Register
		if err := json.Unmarshal(value, &data); err != nil {
			return err
		}

		if info.Password != data.Password {
			return PasswordErrorErr
		}

		id = data.Id

		return nil
	})

	return id, err
}

func (info Register) ChangePassword(newPassword string) error {
	if err := checkPassword(newPassword); err != nil {
		return err
	}
	uname := []byte(info.Username)
	return loginDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(loginBucket)

		value := bucket.Get(uname)
		if value == nil {
			return UserExistErr
		}

		var data Register
		if err := json.Unmarshal(value, &data); err != nil {
			return err
		}

		data.Password = newPassword
		value, err := json.Marshal(data)
		if err != nil {
			return err
		}

		return bucket.Put(uname, value)
	})
}

func IsUserExist(username string) error {
	uname := []byte(username)
	return loginDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(loginBucket)
		value := bucket.Get(uname)
		if value != nil {
			return UserExistErr
		} else {
			return nil
		}
	})
}
