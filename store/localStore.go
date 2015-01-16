package store

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

/*
 * Simple key/value store using Bolt
 */

type (
	LocalClient struct {
		User *User
	}
	Bet struct {
		UserId       string    `json:"source"`
		Feed         string    `json:"feed-address"`
		Deadline     time.Time `json:"deadline"`
		Type         string    `json:"bet-type"` //
		TargetValue  string    `json:"target-value"`
		Wager        float32   `json:"wager"`
		CounterWager float32   `json:"counterwager"`
	}
)

const (
	DATA_COLLECTION = "data"
	USR_COLLECTION  = "user"
	BET_COLLECTION  = "bet"
	storeName       = "intertidal.db"
)

func NewLocalClient(user *User) *LocalClient {

	lc := &LocalClient{User: user}

	if lc.User.IsSet() == false {
		if err := lc.login(); err == nil {
			log.Panicf("No user found: %s", err.Error())
		}
	}

	return lc

}

func saveIt(what interface{}, where string, who string) error {
	jsonData, _ := json.Marshal(what)

	db, err := bolt.Open(storeName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		dataB, err := tx.CreateBucketIfNotExists([]byte(where))
		err = dataB.Put([]byte(who), jsonData)
		return err
	})
	return err
}

func getIt(what, where string, data interface{}) error {
	db, err := bolt.Open(storeName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		dataB, err := tx.CreateBucketIfNotExists([]byte(where))
		jsonData := dataB.Get([]byte(what))

		err = json.Unmarshal(jsonData, &data)
		return err
	})
	return nil
}

// we need to login to the platform to be able to us it
func (lc *LocalClient) login() error {
	return getIt("current", USR_COLLECTION, &lc.User)
}

// we need to login to the platform to be able to us it
func (lc *LocalClient) lodge(b *Bet) error {
	return saveIt(b, BET_COLLECTION, lc.User.Id)
}

// we need to login to the platform to be able to us it
func (lc *LocalClient) check(b *Bet) error {
	return getIt(lc.User.Id, BET_COLLECTION, &b)
}

func (lc *LocalClient) Sync(with Client) error {

	// get what we have locally
	qry := &Query{UserId: lc.User.Id}
	toSync, _ := lc.Run(qry)

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

	return saveIt(data, DATA_COLLECTION, lc.User.Id)
}

func doQuery(all []map[string]interface{}, qry *Query) (matches []interface{}) {

	for i := range all {
		if len(qry.Types) == 0 && qry.FromTime == "" {
			matches = append(matches, all[i])
		} else {

			for t := range qry.Types {

				if qry.FromTime == "" {
					if all[i]["type"] == qry.Types[t] {
						matches = append(matches, all[i])
					}
				} else {
					log.Print("Time not yet implemented")
				}
			}
		}
	}
	return matches
}

func (lc *LocalClient) Run(qry *Query) (results []interface{}, err error) {

	if lc.User.Id == "" {
		return nil, errors.New(USR_ID_NOTSET)
	}

	var data []map[string]interface{}

	if err = getIt(lc.User.Id, DATA_COLLECTION, &data); err != nil {
		return nil, err
	}

	return doQuery(data, qry), err
}
