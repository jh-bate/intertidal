package store

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
)

/*
 * Simple implamentation of key/value store using Bolt
 */

type BoltClient struct{}

const (
	DATA   = "data"
	USER   = "user"
	CONFIG = "config"
	DB     = "intertidal.db"
)

func NewBoltClient() *BoltClient {
	return &BoltClient{}
}

func (b *BoltClient) Ping() error {
	return nil
}
func (b *BoltClient) StoreUser(key, token string) error {

	db, err := bolt.Open(DB, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		dataB, err := tx.CreateBucketIfNotExists([]byte(USER))
		err = dataB.Put([]byte(key), []byte(token))
		return err
	})
	return err
}

func (b *BoltClient) RetrieveUser(key string) (results []interface{}, err error) {

	db, err := bolt.Open(DB, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		dataB, err := tx.CreateBucketIfNotExists([]byte(USER))
		jsonData := dataB.Get([]byte(key))

		err = json.Unmarshal(jsonData, &results)
		return err
	})
	return results, err
}
func (b *BoltClient) StoreUserData(usr string, data []interface{}) error {

	jsonData, _ := json.Marshal(data)

	db, err := bolt.Open(DB, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		dataB, err := tx.CreateBucketIfNotExists([]byte(DATA))
		err = dataB.Put([]byte(usr), jsonData)
		return err
	})
	return err
}
func (b *BoltClient) RetrieveUserData(usr string) (results []interface{}, err error) {

	db, err := bolt.Open(DB, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		dataB, err := tx.CreateBucketIfNotExists([]byte(DATA))
		jsonData := dataB.Get([]byte(usr))

		err = json.Unmarshal(jsonData, &results)
		return err
	})
	return results, err
}

func (b *BoltClient) StoreConfig(key string, data []interface{}) error {
	jsonData, _ := json.Marshal(data)

	db, err := bolt.Open(DB, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		dataB, err := tx.CreateBucketIfNotExists([]byte(CONFIG))
		err = dataB.Put([]byte(key), jsonData)
		return err
	})
	return err
}
func (b *BoltClient) RetrieveConfig(key string) (results []interface{}, err error) {

	db, err := bolt.Open(DB, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		dataB, err := tx.CreateBucketIfNotExists([]byte(CONFIG))
		jsonData := dataB.Get([]byte(key))

		err = json.Unmarshal(jsonData, &results)
		return err
	})
	return results, err
}
