package actions

import (
	"github.com/jh-bate/intertidal/data"
)

type (
	SyncData struct {
		source data.Store
		target data.Store
		name   string
	}
)

func SyncDataAction(calldata interface{}, store data.Store, name string) *SyncData {

	target := calldata.(data.Store)

	return &SyncData{source: store, target: target, name: name}
}

func (a *SyncData) Execute() (interface{}, error) {
	return nil, a.source.Sync(a.name, a.target)
}
