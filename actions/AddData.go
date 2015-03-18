package actions

import (
	"github.com/jh-bate/intertidal/data"
)

type (
	AddData struct {
		data  interface{}
		store data.Store
		name  string
	}
)

func AddDataAction(data interface{}, store data.Store, name string) *AddData {
	return &AddData{data: data, store: store, name: name}
}

func (a *AddData) Execute() error {
	return a.store.Save(a.name, a.data)
}
