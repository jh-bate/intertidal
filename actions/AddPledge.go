package actions

import (
	"github.com/jh-bate/intertidal/data"
)

type (
	AddPledge struct {
		store  data.Store
		pledge *data.Pledge
		name   string
	}
)

func AddPledgeAction(store data.Store, name string, pledge *data.Pledge) *AddPledge {
	return &AddPledge{store: store, name: name, pledge: pledge}
}

func (a *AddPledge) Execute() error {
	return a.store.Save(a.name, a.pledge)
}
