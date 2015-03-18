package actions

import (
	"errors"

	"github.com/jh-bate/intertidal/data"
)

type (
	Manager struct {
		Current Action
	}
	ActionType string
)

const (
	//Available Action's
	ActionTypeAddData   ActionType = "add_data"
	ActionTypeLoadData  ActionType = "load_data"
	ActionTypeQueryData ActionType = "query_data"
	ActionTypeSyncData  ActionType = "sync_data"

	ActionTypeAddPledge   ActionType = "add_pledge"
	ActionTypeCheckPledge ActionType = "check_pledge"
)

func CreateActionManager() *Manager {
	return &Manager{}
}

func (a *Manager) Execute(action ActionType, calldata interface{}, store data.Store, name string) (interface{}, error) {

	switch action {
	case ActionTypeAddData:
		a.Current = AddDataAction(calldata, store, name)
		return a.Current.Execute()
	case ActionTypeLoadData:
		a.Current = LoadDataAction(calldata, store, name)
		return a.Current.Execute()
	case ActionTypeQueryData:
		a.Current = QueryDataAction(calldata, store, name)
		return a.Current.Execute()
	case ActionTypeSyncData:
		a.Current = SyncDataAction(calldata, store, name)
		return a.Current.Execute()
	case ActionTypeAddPledge:
		a.Current = AddPledgeAction(calldata, store, name)
		return a.Current.Execute()
	case ActionTypeCheckPledge:
		a.Current = CheckPledgeAction(calldata, store, name)
		return a.Current.Execute()
	}
	return nil, errors.New("no matching action")
}
