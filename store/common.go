package store

import (
	"fmt"
	"strings"
)

type (
	// Generic store interface
	Client interface {
		Ping() error
		Save(data []interface{}) error
		Query(qry *Query) (data []interface{}, err error)
	}
	// Query
	Query struct {
		UserId   string
		Types    []string
		FromTime string
	}
	// User
	User struct {
		Token string `json:"-"`
		Id    string `json:"-"`
		Name  string `json:"username"`
	}
)

const (
	USR_ID_NOTSET   = "The User.Id is required but hasn't been set"
	USR_NAME_NOTSET = "The User.Name is required but hasn't been set"
)

func (q *Query) ToString() string {
	const queryString = "METAQUERY WHERE userid IS %s QUERY TYPE IN %s WHERE time > %s SORT BY time AS Timestamp REVERSED"
	return fmt.Sprintf(queryString, q.UserId, strings.Join(q.Types, ","), q.FromTime)
}
