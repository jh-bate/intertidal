package actions

import (
	"github.com/jh-bate/intertidal/data"
)

type (
	QueryData struct {
		query   *data.Query
		results interface{}
		store   data.Store
		name    string
	}
)

func QueryDataAction(query *data.Query, store data.Store, name string) *QueryData {
	return &QueryData{query: query, store: store, name: name}
}

func (a *QueryData) Execute() (err error) {
	if a.results, err = a.store.Query2(a.name, a.query); err != nil {
		return err
	}
	return nil
}
