package store

import (
	"fmt"
	"strings"
)

type (
	// Generic store interface
	Store interface {
		Ping() error
		Login() error
		//data
		Find(collection string, results interface{}) error
		Query(qry *Query) (results []map[string]interface{}, err error)
		Save(collection string, data interface{}) error
		Sync(collection string, other Store) error
	}
	// Query
	Query struct {
		UserId   string
		Types    []string
		FromTime string
	}

	QueryOne struct {
		UserId   string
		Type     string
		FromTime string
		Equals   struct {
			Val  float32
			Cond string
		}
	}
)

func (q *Query) ToString() string {
	const queryString = "METAQUERY WHERE userid IS %s QUERY TYPE IN %s WHERE time > %s SORT BY time AS Timestamp REVERSED"
	return fmt.Sprintf(queryString, q.UserId, strings.Join(q.Types, ","), q.FromTime)
}
