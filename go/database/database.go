package database

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

const DBDirectory = "data"
const DBFile = DBDirectory + "/users.db"

var DB *bolt.DB

type User struct {
    WID [16]byte
    HASH [16]byte
}

func init() {
    var err error

    // Make database directory if it doesn't exist
    if _, err := os.Stat(DBDirectory); os.IsNotExist(err) {
	err := os.Mkdir(DBDirectory, 0700)
	if err != nil {
	    log.Fatal(err)
	}
    }

    // Open the database
    DB, err = bolt.Open(DBFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
    if err != nil {
	log.Fatal(err) 
    }


    // Create the list of users if it doesn't exist 
    err = DB.Update(func(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists([]byte("Users"))
	return err
    })
    if err != nil {
	log.Fatal(err)
    }

    // Add testing user
    testUser, err := NewUser("000000000000000000773030376D6171", "55f8b969f2a7c33cfb87edaa2d1afafd")
    if err != nil {
	log.Fatal(err)
    }

    err = InsertUser(testUser)
    if err != nil {
	log.Fatal(err)
    }
}

func NewUser(widHex, hashHex string) (User, error) {
    var u User

    widBytes, err := hex.DecodeString(widHex)
    if err != nil {
	return u, errors.New("Invalid WID")
    }

    if len(widBytes) > 16 {
	return u, errors.New("WID too long")
    }

    hashBytes, err := hex.DecodeString(hashHex)
    if err != nil {
	return u, errors.New("Invalid HASH")
    }

    if len(hashBytes) > 16 {
	return u, errors.New("Hash too long")
    }

    copy(u.WID[:], widBytes)
    copy(u.HASH[:], hashBytes)

    return u, nil
}

func InsertUser(u User) error {
    return DB.Update(func(tx *bolt.Tx) error {
	b := tx.Bucket([]byte("Users"))
	if b == nil {
	    return errors.New("Users bucket not found")
	}

	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, u)
	if err != nil {
	    return err
	}

	return b.Put(u.WID[:], buf.Bytes())
    })
}

func QueryUser(wid [16]byte) (User, error) {
    var u User

    err := DB.View(func(tx *bolt.Tx) error {
	b := tx.Bucket([]byte("Users"))
	if b == nil {
	    return errors.New("Users bucket not found")
	}

	data := b.Get(wid[:])
	if data == nil {
	    return errors.New("User not found")
	}

	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, &u)
    })

    return u, err
}
