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

func AddPledgeAction(calldata interface{}, store data.Store, name string) *AddPledge {

	pledge := calldata.(*data.Pledge)
	return &AddPledge{store: store, name: name, pledge: pledge}
}

func (a *AddPledge) Execute() (interface{}, error) {
	return nil, a.store.Save(a.name, a.pledge)
}
