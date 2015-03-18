package actions

import (
	"github.com/jh-bate/intertidal/data"
)

type (
	QueryData struct {
		query   *data.Query
		Results interface{}
		store   data.Store
		name    string
	}
)

func QueryDataAction(calldata interface{}, store data.Store, name string) *QueryData {

	query := calldata.(*data.Query)
	return &QueryData{query: query, store: store, name: name}
}

func (a *QueryData) Execute() (results interface{}, err error) {
	return a.store.Query(a.name, a.query)
}
