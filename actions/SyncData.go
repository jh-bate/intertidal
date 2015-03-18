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

func SyncDataAction(source data.Store, target data.Store, name string) *SyncData {
	return &SyncData{source: source, target: target, name: name}
}

func (a *SyncData) Execute() error {
	return a.source.Sync(a.name, a.target)
}
