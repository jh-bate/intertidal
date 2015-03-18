package actions

import (
	"errors"

	"github.com/jh-bate/intertidal/data"
)

type (
	ActionManager struct {
		current Action
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

func CreateActionManager() *ActionManager {
	return &ActionManager{}
}

func (a *ActionManager) Execute(action ActionType, calldata interface{}, store data.Store, name string) (interface{}, error) {

	switch action {
	case ActionTypeAddData:
		a.current = AddDataAction(calldata, store, name)
		return a.current.Execute()
	case ActionTypeLoadData:
		a.current = LoadDataAction(calldata, store, name)
		return a.current.Execute()
	case ActionTypeQueryData:
		a.current = QueryDataAction(calldata, store, name)
		return a.current.Execute()
	case ActionTypeSyncData:
		a.current = SyncDataAction(calldata, store, name)
		return a.current.Execute()
	case ActionTypeAddPledge:
		a.current = AddPledgeAction(calldata, store, name)
		return a.current.Execute()
	case ActionTypeCheckPledge:
		a.current = CheckPledgeAction(calldata, store, name)
		return a.current.Execute()
	}
	return nil, errors.New("no matching action")
}
