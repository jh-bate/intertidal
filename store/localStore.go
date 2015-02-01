package store

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/boltdb/bolt"
)

/*
 * Simple key/value store using Bolt
 */

type (
	LocalStore struct {
		User *User
	}
)

const (
	USR_COLLECTION = "user"
	storeName      = "intertidal.db"
)

func NewLocalStore(user *User) *LocalStore {

	lc := &LocalStore{User: user}

	if lc.User.IsSet() == false {
		if err := lc.Login(); err != nil {
			log.Panicf("No user found: %s", err.Error())
		}
	}

	return lc
}

// we need to login to the platform to be able to us it
func (ls *LocalStore) Login() error {
	if err := ls.Find(USR_COLLECTION, ls.User); err != nil {
		return err
	}
	ls.User.Id = "todo2"
	return nil
}

func (ls *LocalStore) Ping() error {
	_, err := bolt.Open(storeName, 0600, nil)
	return err
}

func (ls *LocalStore) Save(collection string, data interface{}) error {
	jsonData, _ := json.Marshal(data)

	db, err := bolt.Open(storeName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		dataB, err := tx.CreateBucketIfNotExists([]byte(collection))
		err = dataB.Put([]byte(ls.User.Id), jsonData)
		return err
	})
	return err
}

func (ls *LocalStore) Find(collection string, results interface{}) error {
	db, err := bolt.Open(storeName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		dataB, err := tx.CreateBucketIfNotExists([]byte(collection))
		jsonData := dataB.Get([]byte(ls.User.Id))

		log.Printf("found this %s", jsonData)

		err = json.Unmarshal(jsonData, &results)
		return err
	})
	return nil
}

func (ls *LocalStore) Query(qry *Query) (results []map[string]interface{}, err error) {
	if ls.User.Id == "" {
		return nil, errors.New(USR_ID_NOTSET)
	}

	var all []map[string]interface{}

	if err = ls.Find(DATA_COLLECTION, &all); err != nil {
		return nil, err
	}

	return doQuery(all, qry), nil
}

func (ls *LocalStore) Sync(collection string, with Store) error {

	// get what we have locally
	qry := &Query{UserId: ls.User.Id}

	var data interface{}
	ls.Find(collection, &data)

	//send it to the other store reporting what occured
	if syncErr := with.Save(collection, data); syncErr != nil {
		log.Printf("Error trying to sync query: %v err: %v ", qry, syncErr)
		return syncErr
	}
	log.Print("Successfully synced")
	return nil
}

func doQuery(all []map[string]interface{}, qry *Query) (results []map[string]interface{}) {

	for i := range all {
		if len(qry.Types) == 0 && qry.FromTime == "" {
			results = append(results, all[i])
		} else {

			for t := range qry.Types {

				if qry.FromTime == "" {
					if all[i]["type"] == qry.Types[t] {
						//log.Printf("match ... %v", all[i])
						results = append(results, all[i])
					}
				} else {
					log.Print("Time not yet implemented")
				}
			}
		}
	}
	return results
}
