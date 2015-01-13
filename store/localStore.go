package store

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/boltdb/bolt"
)

/*
 * Simple implamentation of key/value store using Bolt
 */

type LocalClient struct {
	User *User
}

const (
	DATA_COLLECTION = "data"
	storeName       = "intertidal.db"
)

func NewLocalClient(user *User) *LocalClient {
	return &LocalClient{User: user}
}

func (lc *LocalClient) Sync(with Client) error {

	// get what we have locally
	qry := &Query{UserId: lc.User.Id}
	toSync, _ := lc.Query(qry)

	//send it to the other store reporting what occured
	if syncErr := with.Save(toSync); syncErr != nil {
		log.Printf("Error trying to sync query: %v err: %v ", qry, syncErr)
		return syncErr
	}
	log.Print("Successfully synced")
	return nil
}

func (lc *LocalClient) Ping() error {
	_, err := bolt.Open(storeName, 0600, nil)
	return err
}

func (lc *LocalClient) Save(data []interface{}) error {

	if lc.User.Id == "" {
		return errors.New(USR_ID_NOTSET)
	}

	jsonData, _ := json.Marshal(data)

	db, err := bolt.Open(storeName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		dataB, err := tx.CreateBucketIfNotExists([]byte(DATA_COLLECTION))
		err = dataB.Put([]byte(lc.User.Id), jsonData)
		return err
	})
	return err
}

func (lc *LocalClient) Query(qry *Query) (data []interface{}, err error) {

	if lc.User.Id == "" {
		return nil, errors.New(USR_ID_NOTSET)
	}

	db, err := bolt.Open(storeName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		dataB, err := tx.CreateBucketIfNotExists([]byte(DATA_COLLECTION))
		jsonData := dataB.Get([]byte(lc.User.Id))

		err = json.Unmarshal(jsonData, &data)
		return err
	})
	return data, err
}
