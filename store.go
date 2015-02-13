package intertidal

import (
	"fmt"
	"strings"
)

type (
	Store interface {
		//ping the store to see that it is available
		Ping() error
		//login to this store
		Login() error

		//find all the data in the given collection for the logged in user
		Find(collection string, results interface{}) error
		//find all data the matches the given query and collection
		Query(qry *Query) (results []map[string]interface{}, err error)
		//save the given data to the named collection for the logged in user
		Save(collection string, data interface{}) error
		//the named collection in this store with the given store
		Sync(collection string, other Store) error
	}

	Query struct {
		//run against this users data
		UserId string
		//for these types e.g. smbg, food, note
		Types []string
		//from this given time
		FromTime string
	}
)

func (q *Query) ToString() string {
	const queryString = "METAQUERY WHERE userid IS %s QUERY TYPE IN %s WHERE time > %s SORT BY time AS Timestamp REVERSED"
	return fmt.Sprintf(queryString, q.UserId, strings.Join(q.Types, ","), q.FromTime)
}
